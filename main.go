package main

import (
	"fmt"
	"trends"
)

func main() {
	fmt.Println("🚀 Fetching and Analyzing Bug Bounty Trends...")
	trends.TrackTrends()
	fmt.Println("✅ Trends successfully saved to `trends.md`!")
}
