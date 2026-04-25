# portrank

The `portrank` package assigns a **risk score** to port numbers based on their
exposure profile and historical attack surface.

## Scores

| Constant        | Value | Label    |
|-----------------|-------|----------|
| `ScoreNone`     | 0     | none     |
| `ScoreLow`      | 25    | low      |
| `ScoreMedium`   | 50    | medium   |
| `ScoreHigh`     | 75    | high     |
| `ScoreCritical` | 100   | critical |

## Usage

```go
// Use built-in defaults
r := portrank.New(nil)

// Override or extend defaults
r := portrank.New(map[int]portrank.Score{
    9200: portrank.ScoreHigh, // Elasticsearch
})

score := r.Score(3389)       // 100 (critical) — RDP
fmt.Println(score.Label())   // "critical"
fmt.Println(score.String())  // "100 (critical)"

if r.IsCritical(port) {
    // trigger high-priority alert
}
```

## Built-in Rankings

Well-known risky ports are ranked by default, including FTP (21), Telnet (23),
RDP (3389), SMB (445), and common database ports. Unknown ports receive
`ScoreNone` and do not trigger elevated alerts.
