import os
import requests
import json

def main():
    # GitHub Environment မှ ဒေတာများ ရယူခြင်း
    gemini_key = os.environ.get("GEMINI_API_KEY")
    
    # ပြင်ဆင်လိုက်သော ကုဒ်များကို ဖတ်ခြင်း
    if not os.path.exists("pr_diff.txt"):
        print("No diff file found.")
        return
        
    with open("pr_diff.txt", "r") as f:
        git_diff = f.read()

    if not git_diff.strip():
        print("Diff is empty.")
        return

    # Gemini API သို့ ပေးပို့မည့် Prompt ရေးသားခြင်း
    prompt = f"""
    You are an expert code reviewer. Review the following git diff of a Pull Request.
    Identify any potential bugs, security vulnerabilities, or performance issues.
    Provide constructive feedback and code suggestions in a clear bulleted format.
    Keep the response concise and write it in English.
    
    Git Diff:
    {git_diff}
    """

    # Gemini 1.5 Flash API သို့ Request ပို့ခြင်း
    url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key={gemini_key}"
    headers = {'Content-Type': 'application/json'}
    payload = {
        "contents": [{
            "parts": [{"text": prompt}]
        }]
    }

    response = requests.post(url, headers=headers, data=json.dumps(payload))
    
    if response.status_code == 200:
        result = response.json()
        review_text = result['candidates'][0]['content']['parts'][0]['text']
        
        # ရလဒ်ကို GitHub Markdown Comment အဖြစ် သိမ်းဆည်းခြင်း
        with open("review_result.md", "w") as f:
            f.write("### 🤖 Gemini AI Pull Request Review\n\n")
            f.write(review_text)
    else:
        print(f"API Error: {response.status_code}")
        print(response.text)

if __name__ == "__main__":
    main()