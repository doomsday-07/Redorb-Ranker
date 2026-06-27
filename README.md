# Redrob Hackathon — Intelligent Candidate Ranking System

**Team:** SlEePeR_cElL

Go-based ranker for the India Runs Data & AI Challenge. Ranks 100,000 candidates against a Senior AI Engineer — Founding Team job description at Redrob AI.

## Setup

**Prerequisites:** Go 1.26+ (no third-party dependencies)

```bash
git clone <repo-url>
cd ranker
```

Place `candidates.jsonl` in the project directory.

## Build & Run

```bash
go build -o ranker .
./ranker candidates.jsonl submission.csv
```

Runtime is under 1 minute wall-clock for 100K candidates on CPU (single-threaded; well within the 5-minute limit).

## Output

Generates a CSV with columns:
```
candidate_id,rank,score,reasoning
```

Includes the top 100 candidates ranked best-fit first, with a 1-2 sentence reasoning per candidate. Scores are monotonically non-increasing.

## Validation

```bash
python3 validate_submission.py submission.csv
```

## Scoring Methodology

| Component | Weight | What it measures |
|---|---|---|
| Technical Skill Match | 35% | AI-core vs data-eng vs general vs non-AI skills; career-backing verification; summary keyword signals |
| Career Trajectory & Title | 20% | Title relevance to ML/AI roles; services vs product-company history; AI-role experience ratio |
| Keyword Stuffing Detection | 10% | Penalty for unbacked AI skills, low endorsements, title-skills mismatch |
| Experience Fit | 10% | Sweet spot 5–9 years with peak at 6–8 (per JD) |
| Behavioral Signals | 15% | Recruiter response rate, interview completion, recent activity, profile views, GitHub score |
| Education | 5% | Institution tier (Tier-1 bonus); relevant fields (CS, Data Science, AI, etc.) |
| Location | 5% | India + major tech hub bonus; relocation willingness |
| Honeypot Detection | multiplier | 7 checks: experience inconsistency, expert+zero-endorsements, impossible durations, title-summary mismatch, date inversion, skill count anomalies |

Scoring is purely rule-based with no ML model, no network calls, and no GPU — fully compliant with competition compute constraints.

## Docker Sandbox

Build and run the ranker in an isolated container (simulates Stage 3 reproduction environment):

```bash
docker build -t redrob-ranker .
docker run --rm redrob-ranker
```

The container runs on a 100-candidate sample. The output is written to `/app/submission.csv` inside the container.

## Project Structure

```
├── main.go                    # Ranker implementation
├── go.mod                     # Go module definition
├── Dockerfile                 # Multi-stage Docker build for sandbox reproduction
├── sample_100.jsonl           # 100-line sample for Docker sandbox demo
├── SlEePeR_cElL.csv           # Final submission output
├── SlEePeR_cElL.xlsx          # Excel copy of submission
├── submission_metadata.yaml   # Submission metadata for portal upload
├── .gitignore                 # Excludes compiled binary
└── README.md                  # This file
```
