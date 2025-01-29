import requests
import json
import pandas as pd
from bs4 import BeautifulSoup
from markdownify import markdownify

# APIs & Sources
sources = {
    "HackerOne": "https://hackerone.com/hacktivity.json",
    "Bugcrowd": "https://bugcrowd.com/programs.json",
    "Intigriti": "https://www.intigriti.com/researcher/bugbounty-programs",
}

def fetch_h1_reports():
    """Fetch latest reports from HackerOne"""
    response = requests.get(sources["HackerOne"])
    if response.status_code == 200:
        try:
            return response.json().get("reports", [])
        except json.JSONDecodeError:
            print("Error: Received invalid JSON from HackerOne API")
            return []
    return []

def fetch_bugcrowd_programs():
    """Fetch active programs from Bugcrowd"""
    response = requests.get(sources["Bugcrowd"])
    if response.status_code == 200:
        return response.json()
    return []

def fetch_intigriti_programs():
    """Fetch active programs from Intigriti"""
    response = requests.get(sources["Intigriti"])
    soup = BeautifulSoup(response.text, "html.parser")
    programs = [a.text.strip() for a in soup.select(".program-card__title")]
    return programs

def extract_trends():
    """Analyze trends from fresh bug bounty reports"""
    reports = fetch_h1_reports()
    keywords = {}
    
    for report in reports:
        title = report["title"].lower()
        for word in title.split():
            keywords[word] = keywords.get(word, 0) + 1

    trending = sorted(keywords.items(), key=lambda x: x[1], reverse=True)[:10]
    
    return trending

def generate_markdown_report(trending):
    """Create a Markdown report from the trending vulnerabilities"""
    markdown = "# 🔥 Trending Bug Bounty Vulnerabilities\n\n"
    markdown += "**Latest insights from bug bounty reports:**\n\n"

    for term, count in trending:
        markdown += f"- **{term}** ({count} reports)\n"

    markdown += "\n\n📌 Stay ahead of new vulnerabilities! 🚀"
    
    with open("reports/latest_trends.md", "w") as file:
        file.write(markdown)

if __name__ == "__main__":
    trending_vulnerabilities = extract_trends()
    generate_markdown_report(trending_vulnerabilities)
