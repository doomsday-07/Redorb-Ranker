package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

// ======================== DATA STRUCTURES ========================

type Candidate struct {
	ID           string   `json:"candidate_id"`
	Profile      Profile  `json:"profile"`
	CareerHistory []Career `json:"career_history"`
	Education    []Edu    `json:"education"`
	Skills       []Skill  `json:"skills"`
	RedrobSignals Signals `json:"redrob_signals"`
}

type Profile struct {
	AnonymizedName string  `json:"anonymized_name"`
	Headline       string  `json:"headline"`
	Summary        string  `json:"summary"`
	Location       string  `json:"location"`
	Country        string  `json:"country"`
	YOE            float64 `json:"years_of_experience"`
	CurrentTitle   string  `json:"current_title"`
	CurrentCompany string  `json:"current_company"`
	CompanySize    string  `json:"current_company_size"`
	Industry       string  `json:"current_industry"`
}

type Career struct {
	Company      string  `json:"company"`
	Title        string  `json:"title"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date"`
	Duration     int     `json:"duration_months"`
	IsCurrent    bool    `json:"is_current"`
	Industry     string  `json:"industry"`
	CompanySize  string  `json:"company_size"`
	Description  string  `json:"description"`
}

type Edu struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field_of_study"`
	StartYear   int    `json:"start_year"`
	EndYear     int    `json:"end_year"`
	Grade       string `json:"grade"`
	Tier        string `json:"tier"`
}

type Skill struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
	Endorsements int   `json:"endorsements"`
	Duration    int    `json:"duration_months"`
}

type SalaryRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Signals struct {
	ProfileCompleteness    float64        `json:"profile_completeness_score"`
	SignupDate             string         `json:"signup_date"`
	LastActiveDate         string         `json:"last_active_date"`
	OpenToWork             bool           `json:"open_to_work_flag"`
	ProfileViews30d        int            `json:"profile_views_received_30d"`
	Applications30d        int            `json:"applications_submitted_30d"`
	RecruiterResponseRate  float64        `json:"recruiter_response_rate"`
	AvgResponseTimeHrs     float64        `json:"avg_response_time_hours"`
	SkillAssessments       map[string]float64 `json:"skill_assessment_scores"`
	ConnectionCount        int            `json:"connection_count"`
	EndorsementsReceived   int            `json:"endorsements_received"`
	NoticePeriodDays       int            `json:"notice_period_days"`
	SalaryRange            SalaryRange    `json:"expected_salary_range_inr_lpa"`
	PreferredWorkMode      string         `json:"preferred_work_mode"`
	WillingToRelocate      bool           `json:"willing_to_relocate"`
	GithubActivity         float64        `json:"github_activity_score"`
	SearchAppearance30d    int            `json:"search_appearance_30d"`
	SavedByRecruiters30d   int            `json:"saved_by_recruiters_30d"`
	InterviewCompletionRate float64       `json:"interview_completion_rate"`
	OfferAcceptanceRate    float64        `json:"offer_acceptance_rate"`
	VerifiedEmail          bool           `json:"verified_email"`
	VerifiedPhone          bool           `json:"verified_phone"`
	LinkedinConnected      bool           `json:"linkedin_connected"`
}

// ======================== SKILL TAXONOMY ========================

// AI/ML Core Skills — high value for this JD
var aiCoreSkills = map[string]bool{
	"machine learning": true, "deep learning": true, "nlp": true,
	"natural language processing": true, "computer vision": true,
	"reinforcement learning": true, "pytorch": true, "tensorflow": true,
	"keras": true, "scikit-learn": true, "sklearn": true, "xgboost": true,
	"lightgbm": true, "catboost": true, "neural networks": true,
	"transformers": true, "hugging face": true, "huggingface": true,
	"bert": true, "gpt": true, "llm": true, "large language models": true,
	"fine-tuning": true, "fine tuning": true, "fine-tuning llms": true,
	"fine-tuning language models": true, "lora": true, "qlora": true,
	"peft": true, "rlhf": true, "prompt engineering": true,
	"rag": true, "retrieval augmented generation": true,
	"vector search": true, "vector databases": true,
	"embeddings": true, "word embeddings": true, "sentence embeddings": true,
	"semantic search": true, "information retrieval": true,
	"milvus": true, "pinecone": true, "weaviate": true, "qdrant": true,
	"faiss": true, "chroma": true, "annoy": true, "hnsw": true,
	"elasticsearch": true, "opensearch": true,
	"ranknet": true, "lambdarank": true, "listwise": true,
	"pairwise": true, "learning to rank": true,
	"recommendation systems": true, "recommender systems": true,
	"text classification": true, "sentiment analysis": true,
	"named entity recognition": true, "ner": true,
	"speech recognition": true, "asr": true, "tts": true,
	"text to speech": true, "speech to text": true,
	"image classification": true, "object detection": true,
	"image segmentation": true, "ocr": true,
	"generative ai": true, "gen ai": true, "genai": true,
	"langchain": true, "llamaindex": true, "llama index": true,
	"openai": true, "anthropic": true, "claude": true,
	"data science": true, "statistical modeling": true,
	"feature engineering": true, "a/b testing": true,
	"experiment design": true, "causal inference": true,
	"bayesian statistics": true, "time series": true,
	"anomaly detection": true, "clustering": true,
	"dimensionality reduction": true, "pca": true,
	"mlops": true, "ml pipeline": true, "model deployment": true,
	"model serving": true, "bentoml": true, "mlflow": true,
	"weights & biases": true, "wandb": true, "dvc": true,
	"kubeflow": true, "seldon": true, "bento": true,
	"onnx": true, "tensorrt": true, "openvino": true,
}

// Data Engineering Skills — supporting value for this JD
var dataEngSkills = map[string]bool{
	"python": true, "sql": true, "spark": true, "pyspark": true,
	"airflow": true, "dbt": true, "snowflake": true, "bigquery": true,
	"redshift": true, "databricks": true, "kafka": true,
	"data pipelines": true, "etl": true, "data warehousing": true,
	"hadoop": true, "hive": true, "presto": true, "trino": true,
	"aws": true, "gcp": true, "azure": true,
	"docker": true, "kubernetes": true, "k8s": true,
	"terraform": true, "ci/cd": true, "git": true,
	"scala": true, "java": true, "go": true, "golang": true,
	"rust": true, "c++": true, "r": true, "julia": true,
}

// Software/Backend Skills — lower value
var softwareSkills = map[string]bool{
	"javascript": true, "typescript": true, "react": true, "angular": true,
	"vue": true, "node.js": true, "nodejs": true, "express": true,
	"django": true, "flask": true, "fastapi": true, "rest api": true,
	"graphql": true, "microservices": true, "redis": true,
	"postgresql": true, "mysql": true, "mongodb": true,
	"html": true, "css": true, "sass": true, "webpack": true,
	"redux": true, "tailwind": true, "bootstrap": true,
	"java": true, "spring": true, ".net": true, "c#": true,
	"php": true, "ruby": true, "rails": true, "swift": true,
	"kotlin": true, "android": true, "ios": true, "react native": true,
	"flutter": true, "dart": true,
}

// Non-AI Specialization Skills — should NOT dominate
var nonAISpecSkills = map[string]bool{
	"photoshop": true, "illustrator": true, "figma": true,
	"graphic design": true, "ui/ux": true, "ux design": true,
	"video editing": true, "animation": true, "3d modeling": true,
	"blender": true, "after effects": true, "premiere pro": true,
	"accounting": true, "finance": true, "financial modeling": true,
	"excel": true, "powerpoint": true, "tableau": true, "power bi": true,
	"seo": true, "content writing": true, "marketing": true,
	"sales": true, "crm": true, "salesforce": true,
	"project management": true, "agile": true, "scrum": true,
	"six sigma": true, "pmp": true, "itil": true,
	"sap": true, "oracle": true, "erp": true,
	"customer support": true, "customer service": true,
	"hr management": true, "recruiting": true, "talent acquisition": true,
	"mechanical engineering": true, "civil engineering": true,
	"electrical engineering": true, "cad": true, "solidworks": true,
	"autocad": true, "ansys": true, "fea": true, "cfd": true,
}

// Services companies — penalty per JD
var servicesCompanies = map[string]bool{
	"tcs": true, "infosys": true, "wipro": true, "accenture": true,
	"cognizant": true, "capgemini": true, "tech mahindra": true,
	"hcl": true, "mindtree": true, "ltimindtree": true,
	"mphasis": true, "hexaware": true, "persistent systems": true,
	"zoho": true, "freshworks": true, "niit": true,
	"bsnl": true, "jio": true,
}

// Product companies — bonus
var productCompanies = map[string]bool{
	"google": true, "meta": true, "facebook": true, "amazon": true,
	"apple": true, "microsoft": true, "netflix": true, "uber": true,
	"airbnb": true, "stripe": true, "spotify": true, "twitter": true,
	"snowflake": true, "databricks": true, "openai": true, "anthropic": true,
	"huggingface": true, "palantir": true, "scale ai": true,
	"instacart": true, "doordash": true, "lyft": true, "shopify": true,
	"salesforce": true, "adobe": true, "vmware": true,
	"razorpay": true, "flipkart": true, "swiggy": true, "zomato": true,
	"phonepe": true, "cred": true, "meesho": true, "byju": true,
	"ola": true, "paytm": true, "freshworks": true, "zoho": true,
	"atlassian": true, "twilio": true, "cloudflare": true,
	"pied piper": true, "hooli": true, "stark industries": true,
	"wayne enterprises": true, "globex inc": true, "acme corp": true,
	"initech": true, "dunder mifflin": true,
}

// ======================== SCORING ========================

func scoreCandidate(c *Candidate) (float64, string) {
	totalScore := 0.0
	reasons := []string{}

	// 1. Technical Skill Match (0.35 weight)
	techScore, techReason := scoreTechnicalSkills(c)
	totalScore += techScore * 0.35
	if techReason != "" {
		reasons = append(reasons, techReason)
	}

	// 2. Career Trajectory & Title Fit (0.20 weight)
	careerScore, careerReason := scoreCareerTrajectory(c)
	totalScore += careerScore * 0.20
	if careerReason != "" {
		reasons = append(reasons, careerReason)
	}

	// 3. Keyword Stuffing Detection (0.10 weight — penalty)
	stuffScore, stuffReason := scoreKeywordStuffing(c)
	totalScore += stuffScore * 0.10
	if stuffReason != "" {
		reasons = append(reasons, stuffReason)
	}

	// 4. Experience Fit (0.10 weight)
	expScore, expReason := scoreExperienceFit(c)
	totalScore += expScore * 0.10
	if expReason != "" {
		reasons = append(reasons, expReason)
	}

	// 5. Behavioral Signals (0.15 weight)
	behavScore, behavReason := scoreBehavioralSignals(c)
	totalScore += behavScore * 0.15
	if behavReason != "" {
		reasons = append(reasons, behavReason)
	}

	// 6. Education (0.05 weight)
	eduScore, eduReason := scoreEducation(c)
	totalScore += eduScore * 0.05
	if eduReason != "" {
		reasons = append(reasons, eduReason)
	}

	// 7. Location (0.05 weight)
	locScore, locReason := scoreLocation(c)
	totalScore += locScore * 0.05
	if locReason != "" {
		reasons = append(reasons, locReason)
	}

	// Honeypot penalty
	honeypotPenalty := honeypotScore(c)
	totalScore *= honeypotPenalty

	return totalScore, strings.Join(reasons, "; ")
}

func proficiencyWeight(p string) float64 {
	switch strings.ToLower(p) {
	case "expert":
		return 1.0
	case "advanced":
		return 0.85
	case "intermediate":
		return 0.6
	case "beginner":
		return 0.3
	default:
		return 0.3
	}
}

func normalizeName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func skillInDescription(skillName, description string) bool {
	name := normalizeName(skillName)
	desc := normalizeName(description)

	// Direct match
	if strings.Contains(desc, name) {
		return true
	}

	// Abbreviation/alias matches
	aliases := map[string][]string{
		"machine learning":       {"ml ", " ml"},
		"deep learning":          {"dl ", " dl"},
		"natural language processing": {"nlp"},
		"large language models":  {"llm"},
		"fine-tuning llms":       {"fine-tun", "finetun"},
		"retrieval augmented generation": {"rag"},
		"computer vision":        {"cv ", " cv"},
		"reinforcement learning":  {"rl ", " rl"},
		"speech recognition":     {"asr", "speech"},
		"text to speech":         {"tts"},
		"image classification":   {"image class", "image recog"},
		"object detection":       {"object det", "detect"},
		"named entity recognition": {"ner"},
		"sentiment analysis":     {"sentiment"},
		"recommendation systems": {"recomm", "recommend"},
		"vector search":          {"vector", "embed", "similarity"},
		"information retrieval":  {"retriev", "search", "rank"},
		"learning to rank":       {"ltr", "rank"},
		"data science":           {"data sci"},
		"feature engineering":    {"feature eng", "feature build"},
		"a/b testing":            {"ab test", "a/b test", "experiment"},
		"mlops":                  {"mlops", "ml ops", "ml pipeline"},
	}

	if aliases, ok := aliases[name]; ok {
		for _, a := range aliases {
			if strings.Contains(desc, a) {
				return true
			}
		}
	}

	return false
}

func scoreTechnicalSkills(c *Candidate) (float64, string) {
	if len(c.Skills) == 0 {
		return 0, ""
	}

	aiCoreCount := 0
	dataEngCount := 0
	softwareCount := 0
	nonAICount := 0
	aiSkillScore := 0.0
	totalSkills := len(c.Skills)

	// Combine all descriptions for career-level matching
	allDesc := ""
	for _, ch := range c.CareerHistory {
		allDesc += " " + ch.Description
	}
	allDesc += " " + c.Profile.Summary

	for _, sk := range c.Skills {
		name := normalizeName(sk.Name)
		pwt := proficiencyWeight(sk.Proficiency)
		endorseBonus := math.Min(float64(sk.Endorsements)/30.0, 1.0) * 0.2

		if aiCoreSkills[name] {
			aiCoreCount++
			skillScore := pwt + endorseBonus
			// Bonus if skill is actually used in career (not just listed)
			if skillInDescription(sk.Name, allDesc) {
				skillScore += 0.15
			}
			// Bonus for advanced/expert proficiency in AI skills
			if sk.Proficiency == "advanced" || sk.Proficiency == "expert" {
				skillScore += 0.1
			}
			aiSkillScore += skillScore
		} else if dataEngSkills[name] {
			dataEngCount++
			aiSkillScore += (pwt + endorseBonus) * 0.4 // lower weight
		} else if softwareSkills[name] {
			softwareCount++
			aiSkillScore += (pwt + endorseBonus) * 0.2
		} else if nonAISpecSkills[name] {
			nonAICount++
			aiSkillScore += (pwt + endorseBonus) * 0.05
		}
	}

	// Normalize to 0-1
	if totalSkills > 0 {
		aiSkillScore /= float64(totalSkills)
	}

	// Bonus for AI concentration (more AI skills = more relevant)
	aiRatio := 0.0
	if totalSkills > 0 {
		aiRatio = float64(aiCoreCount) / float64(totalSkills)
	}
	aiSkillScore += aiRatio * 0.3

	// Penalty if >60% skills are non-AI specialty
	if totalSkills > 0 && float64(nonAICount)/float64(totalSkills) > 0.6 {
		aiSkillScore *= 0.5
	}

	// Summary-level signal: check for AI/ML keywords in summary
	summaryScore := 0.0
	lowerSummary := strings.ToLower(c.Profile.Summary)
	summaryKeywords := []string{
		"ml engineer", "machine learning", "ai engineer", "artificial intelligence",
		"data scientist", "deep learning", "nlp", "recommendation",
		"ranking", "retrieval", "embeddings", "llm", "generative ai",
		"production ml", "applied ml", "ml systems", "ai systems",
		"search systems", "information retrieval",
	}
	for _, kw := range summaryKeywords {
		if strings.Contains(lowerSummary, kw) {
			summaryScore += 0.05
		}
	}
	summaryScore = math.Min(summaryScore, 0.3)

	reason := fmt.Sprintf("%d/%d AI-core skills", aiCoreCount, totalSkills)
	return math.Min(aiSkillScore+summaryScore, 1.0), reason
}

func scoreCareerTrajectory(c *Candidate) (float64, string) {
	score := 0.5 // base

	lowerTitle := strings.ToLower(c.Profile.CurrentTitle)
	lowerSummary := strings.ToLower(c.Profile.Summary)

	// Strong signals from title
	titleBonuses := map[string]float64{
		"ml engineer": 0.35, "machine learning engineer": 0.35,
		"ai engineer": 0.35, "artificial intelligence engineer": 0.35,
		"data scientist": 0.25, "data science": 0.25,
		"research scientist": 0.15, "research engineer": 0.2,
		"sr. data scientist": 0.3, "senior data scientist": 0.3,
		"nlp engineer": 0.35, "speech engineer": 0.2,
		"computer vision engineer": 0.15, "cv engineer": 0.15,
		"recommendation engineer": 0.3, "search engineer": 0.3,
		"retrieval engineer": 0.3, "ranking engineer": 0.3,
		"mlops engineer": 0.25, "ml platform engineer": 0.25,
		"ai platform engineer": 0.25, "ai infrastructure engineer": 0.25,
		"data engineer": 0.15, "analytics engineer": 0.1,
		"software engineer": 0.05, "backend engineer": 0.05,
		"full stack engineer": 0.05, "fullstack engineer": 0.05,
		"frontend engineer": 0.0, "mobile developer": 0.0,
		"devops engineer": 0.05, "sre": 0.05,
		"product manager": 0.0, "project manager": 0.0,
		"marketing manager": -0.15, "operations manager": -0.1,
		"accountant": -0.1, "hr manager": -0.1,
		"customer support": -0.15, "mechanical engineer": -0.1,
		"civil engineer": -0.1, "electrical engineer": -0.05,
		"graphic designer": -0.1, "content writer": -0.1,
		"sales": -0.1, "business analyst": -0.05,
	}

	bestTitleBonus := 0.0
	for pattern, bonus := range titleBonuses {
		if strings.Contains(lowerTitle, pattern) {
			if bonus > bestTitleBonus {
				bestTitleBonus = bonus
			}
		}
	}
	score += bestTitleBonus

	// Check summary for AI/ML role signals
	summaryTitleSignals := []string{
		"ml engineer", "machine learning", "ai engineer", "data scientist",
		"nlp", "deep learning", "recommendation", "ranking", "retrieval",
		"embeddings", "search engineer",
	}
	for _, sig := range summaryTitleSignals {
		if strings.Contains(lowerSummary, sig) {
			score += 0.05
			break
		}
	}

	// Services company penalty
	servicesCount := 0
	productCount := 0
	for _, ch := range c.CareerHistory {
		compLower := normalizeName(ch.Company)
		if servicesCompanies[compLower] {
			servicesCount++
		}
		if productCompanies[compLower] {
			productCount++
		}
	}
	totalJobs := len(c.CareerHistory)
	if totalJobs > 0 {
		servicesRatio := float64(servicesCount) / float64(totalJobs)
		if servicesRatio >= 1.0 {
			score -= 0.15
		} else if servicesRatio >= 0.5 {
			score -= 0.08
		}
		if productCount > 0 {
			score += 0.1
		}
	}

	// Experience in AI/ML roles (check career descriptions for AI keywords)
	aiRoleMonths := 0
	for _, ch := range c.CareerHistory {
		descLower := strings.ToLower(ch.Description)
		titleLower := strings.ToLower(ch.Title)
		isAIRole := false
		aiRoleKeywords := []string{
			"ml", "machine learning", "ai ", "artificial intelligence",
			"nlp", "deep learning", "data science", "recommendation",
			"ranking", "retrieval", "embedding", "llm", "neural",
			"model", "prediction", "classification", "regression",
		}
		for _, kw := range aiRoleKeywords {
			if strings.Contains(titleLower, kw) || strings.Contains(descLower, kw) {
				isAIRole = true
				break
			}
		}
		if isAIRole {
			aiRoleMonths += ch.Duration
		}
	}
	yoe := c.Profile.YOE
	if yoe > 0 {
		aiRatio := float64(aiRoleMonths) / (yoe * 12)
		if aiRatio > 0.5 {
			score += 0.15
		} else if aiRatio > 0.3 {
			score += 0.08
		} else if aiRatio > 0.1 {
			score += 0.03
		}
	}

	reason := ""
	if bestTitleBonus > 0.2 {
		reason = fmt.Sprintf("AI-relevant title (%s)", c.Profile.CurrentTitle)
	} else if bestTitleBonus < 0 {
		reason = fmt.Sprintf("Non-AI title (%s)", c.Profile.CurrentTitle)
	}

	return math.Max(0, math.Min(score, 1.0)), reason
}

func scoreKeywordStuffing(c *Candidate) (float64, string) {
	// High AI skill count with low endorsements + no career evidence = stuffing
	aiSkillCount := 0
	totalEndorsements := 0
	lowEndorsementCount := 0

	for _, sk := range c.Skills {
		name := normalizeName(sk.Name)
		if aiCoreSkills[name] {
			aiSkillCount++
			totalEndorsements += sk.Endorsements
			if sk.Endorsements < 5 {
				lowEndorsementCount++
			}
		}
	}

	// Check if AI skills are backed by career history
	allDesc := ""
	for _, ch := range c.CareerHistory {
		allDesc += " " + ch.Description
	}

	backedSkills := 0
	for _, sk := range c.Skills {
		name := normalizeName(sk.Name)
		if aiCoreSkills[name] && skillInDescription(sk.Name, allDesc) {
			backedSkills++
		}
	}

	score := 1.0

	// Penalty: many AI skills but low endorsements
	if aiSkillCount > 5 && float64(lowEndorsementCount)/float64(aiSkillCount) > 0.6 {
		score -= 0.3
	}

	// Penalty: AI skills not backed by career history
	if aiSkillCount > 3 && backedSkills == 0 {
		score -= 0.4
	} else if aiSkillCount > 3 {
		backingRatio := float64(backedSkills) / float64(aiSkillCount)
		if backingRatio < 0.3 {
			score -= 0.2
		}
	}

	// Penalty: high AI skill count but non-technical role
	lowerTitle := strings.ToLower(c.Profile.CurrentTitle)
	nonTechTitles := []string{"marketing manager", "accountant", "hr manager",
		"operations manager", "customer support", "content writer",
		"graphic designer", "mechanical engineer", "civil engineer",
		"business analyst", "sales"}
	for _, nt := range nonTechTitles {
		if strings.Contains(lowerTitle, nt) && aiSkillCount > 4 {
			score -= 0.3
			break
		}
	}

	reason := ""
	if score < 0.8 {
		reason = fmt.Sprintf("keyword stuffing risk (%d AI skills, %d backed)", aiSkillCount, backedSkills)
	}

	return math.Max(0, math.Min(score, 1.0)), reason
}

func scoreExperienceFit(c *Candidate) (float64, string) {
	yoe := c.Profile.YOE

	// Ideal range: 5-9 years, sweet spot 6-8
	var score float64
	if yoe >= 6 && yoe <= 8 {
		score = 1.0
	} else if yoe >= 5 && yoe <= 9 {
		score = 0.85
	} else if yoe >= 4 && yoe <= 10 {
		score = 0.65
	} else if yoe >= 3 && yoe <= 12 {
		score = 0.45
	} else if yoe >= 2 && yoe <= 15 {
		score = 0.3
	} else {
		score = 0.15
	}

	reason := fmt.Sprintf("%.1f yrs YOE", yoe)
	return score, reason
}

func scoreBehavioralSignals(c *Candidate) (float64, string) {
	s := &c.RedrobSignals
	score := 0.0

	// Recruiter response rate (0-0.25)
	score += s.RecruiterResponseRate * 0.25

	// Interview completion rate (0-0.15)
	score += s.InterviewCompletionRate * 0.15

	// Open to work flag (0-0.12)
	if s.OpenToWork {
		score += 0.12
	}

	// Profile views (normalized, 0-0.10)
	viewsScore := math.Min(float64(s.ProfileViews30d)/50.0, 1.0)
	score += viewsScore * 0.10

	// Saved by recruiters (normalized, 0-0.10)
	savedScore := math.Min(float64(s.SavedByRecruiters30d)/10.0, 1.0)
	score += savedScore * 0.10

	// Recent activity (0-0.10)
	if s.LastActiveDate != "" {
		lastActive, err := time.Parse("2006-01-02", s.LastActiveDate)
		if err == nil {
			daysSince := time.Since(lastActive).Hours() / 24
			if daysSince <= 30 {
				score += 0.10
			} else if daysSince <= 90 {
				score += 0.07
			} else if daysSince <= 180 {
				score += 0.04
			} else if daysSince <= 365 {
				score += 0.02
			}
		}
	}

	// Profile completeness (0-0.08)
	score += (s.ProfileCompleteness / 100.0) * 0.08

	// GitHub activity (0-0.05)
	if s.GithubActivity >= 0 {
		score += math.Min(s.GithubActivity/100.0, 1.0) * 0.05
	}

	// Notice period — shorter is better (0-0.05)
	if s.NoticePeriodDays <= 30 {
		score += 0.05
	} else if s.NoticePeriodDays <= 60 {
		score += 0.03
	} else if s.NoticePeriodDays <= 90 {
		score += 0.01
	}

	// Verification (0-0.03)
	if s.VerifiedEmail {
		score += 0.015
	}
	if s.VerifiedPhone {
		score += 0.015
	}

	reason := fmt.Sprintf("response rate %.0f%%, interview %.0f%%",
		s.RecruiterResponseRate*100, s.InterviewCompletionRate*100)
	if s.OpenToWork {
		reason += ", open to work"
	}

	return math.Min(score, 1.0), reason
}

func scoreEducation(c *Candidate) (float64, string) {
	if len(c.Education) == 0 {
		return 0.3, ""
	}

	score := 0.3 // base for having education
	tier1Count := 0
	tier2Count := 0
	relevantCount := 0

	for _, edu := range c.Education {
		switch edu.Tier {
		case "tier_1":
			tier1Count++
			score += 0.25
		case "tier_2":
			tier2Count++
			score += 0.15
		case "tier_3":
			score += 0.05
		}

		fieldLower := strings.ToLower(edu.Field)
		relevantFields := []string{
			"computer science", "computer engineering", "information technology",
			"data science", "artificial intelligence", "machine learning",
			"mathematics", "statistics", "electronics", "electrical",
			"physics", "computing", "software",
		}
		for _, rf := range relevantFields {
			if strings.Contains(fieldLower, rf) {
				relevantCount++
				score += 0.1
				break
			}
		}
	}

	reason := ""
	if tier1Count > 0 {
		reason = "tier-1 institution"
	} else if tier2Count > 0 {
		reason = "tier-2 institution"
	}

	return math.Min(score, 1.0), reason
}

func scoreLocation(c *Candidate) (float64, string) {
	score := 0.5
	lowerCountry := strings.ToLower(c.Profile.Country)
	lowerLocation := strings.ToLower(c.Profile.Location)

	indiaCities := []string{
		"pune", "noida", "bangalore", "bengaluru", "hyderabad",
		"mumbai", "delhi", "gurgaon", "gurugram", "chennai",
		"kolkata", "ahmedabad", "jaipur", "lucknow", "trivandrum",
		"trivandrum", "kochi", "coimbatore", "indore", "bhopal",
		"chandigarh", "vizag", "bhubaneswar", "vadodara", "surat",
	}

	if lowerCountry == "india" {
		score += 0.2
		for _, city := range indiaCities {
			if strings.Contains(lowerLocation, city) {
				score += 0.15
				break
			}
		}
	} else if lowerCountry == "usa" || lowerCountry == "uk" || lowerCountry == "canada" || lowerCountry == "australia" {
		score += 0.05
	}

	if c.RedrobSignals.WillingToRelocate {
		score += 0.1
	}

	reason := ""
	if lowerCountry == "india" {
		reason = fmt.Sprintf("India (%s)", c.Profile.Location)
	} else {
		reason = fmt.Sprintf("%s", c.Profile.Country)
	}

	return math.Min(score, 1.0), reason
}

func honeypotScore(c *Candidate) float64 {
	penalty := 1.0

	// Check 1: Experience inconsistency
	// If YOE claims much more than sum of career durations
	totalCareerMonths := 0
	for _, ch := range c.CareerHistory {
		totalCareerMonths += ch.Duration
	}
	careerYears := float64(totalCareerMonths) / 12.0
	if c.Profile.YOE > careerYears+3 {
		penalty *= 0.3
	}

	// Check 2: Too many expert skills with 0 endorsements
	expertNoEndorse := 0
	for _, sk := range c.Skills {
		if sk.Proficiency == "expert" && sk.Endorsements == 0 {
			expertNoEndorse++
		}
	}
	if expertNoEndorse >= 3 {
		penalty *= 0.4
	}

	// Check 3: Impossible skill durations (>24 months for a beginner with 1 yr experience)
	if c.Profile.YOE < 2 {
		for _, sk := range c.Skills {
			if sk.Duration > int(c.Profile.YOE*12) && sk.Duration > 24 {
				penalty *= 0.5
				break
			}
		}
	}

	// Check 4: Summary/title mismatch (AI skills listed but summary is pure marketing)
	lowerSummary := strings.ToLower(c.Profile.Summary)
	lowerTitle := strings.ToLower(c.Profile.CurrentTitle)
	aiSummary := strings.Contains(lowerSummary, "ml") ||
		strings.Contains(lowerSummary, "machine learning") ||
		strings.Contains(lowerSummary, "data science") ||
		strings.Contains(lowerSummary, "ai engineer")
	nonAITitle := strings.Contains(lowerTitle, "marketing") ||
		strings.Contains(lowerTitle, "accountant") ||
		strings.Contains(lowerTitle, "hr ") ||
		strings.Contains(lowerTitle, "operations manager") ||
		strings.Contains(lowerTitle, "customer support") ||
		strings.Contains(lowerTitle, "mechanical") ||
		strings.Contains(lowerTitle, "civil engineer")

	aiSkillCount := 0
	for _, sk := range c.Skills {
		if aiCoreSkills[normalizeName(sk.Name)] {
			aiSkillCount++
		}
	}
	if aiSkillCount > 5 && !aiSummary && nonAITitle {
		penalty *= 0.5
	}

	// Check 5: Education year inconsistencies
	for _, edu := range c.Education {
		if edu.EndYear < edu.StartYear {
			penalty *= 0.3
		}
		if edu.EndYear > 2030 {
			penalty *= 0.3
		}
	}

	// Check 6: Career history date inconsistencies
	for _, ch := range c.CareerHistory {
		startDate, err1 := time.Parse("2006-01-02", ch.StartDate)
		if err1 != nil {
			continue
		}
		if ch.EndDate != nil {
			endDate, err2 := time.Parse("2006-01-02", *ch.EndDate)
			if err2 == nil && endDate.Before(startDate) {
				penalty *= 0.3
			}
		}
	}

	// Check 7: Very high skill count (30+) with all expert
	if len(c.Skills) > 30 {
		allExpert := true
		for _, sk := range c.Skills {
			if sk.Proficiency != "expert" && sk.Proficiency != "advanced" {
				allExpert = false
				break
			}
		}
		if allExpert {
			penalty *= 0.4
		}
	}

	return penalty
}

// ======================== MAIN ========================

func main() {
	start := time.Now()
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <candidates.jsonl> <output.csv>")
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", inputPath, err)
	}
	defer file.Close()

	var candidates []Candidate
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var c Candidate
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			fmt.Fprintf(os.Stderr, "Line %d: parse error: %v\n", lineNum, err)
			continue
		}
		candidates = append(candidates, c)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Read error: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Loaded %d candidates\n", len(candidates))

	type Ranked struct {
		CandidateID string
		Score       float64
		Reasoning   string
	}

	ranked := make([]Ranked, len(candidates))
	for i := range candidates {
		score, reasoning := scoreCandidate(&candidates[i])
		ranked[i] = Ranked{
			CandidateID: candidates[i].ID,
			Score:       score,
			Reasoning:   reasoning,
		}
	}

	// Sort by score descending, tie-break by candidate_id ascending
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score != ranked[j].Score {
			return ranked[i].Score > ranked[j].Score
		}
		return ranked[i].CandidateID < ranked[j].CandidateID
	})

	// Take top 100
	top100 := ranked[:100]
	if len(ranked) < 100 {
		top100 = ranked
	}

	// Write CSV
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", outputPath, err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Header
	writer.Write([]string{"candidate_id", "rank", "score", "reasoning"})

	// Enforce strict non-increasing scores with proper tie-breaking
	// Assign final scores ensuring monotonicity
	for i, r := range top100 {
		scoreStr := fmt.Sprintf("%.6f", r.Score)
		writer.Write([]string{r.CandidateID, fmt.Sprintf("%d", i+1), scoreStr, r.Reasoning})
	}

	fmt.Fprintf(os.Stderr, "Written %d candidates to %s\n", len(top100), outputPath)
	fmt.Fprintf(os.Stderr, "Total time: %v\n", time.Since(start))
}
