from flask import Flask, request, jsonify
import openai
import os

app = Flask(__name__)
openai.api_key = os.getenv("OPENAI_API_KEY")

@app.route('/analyze', methods=['POST'])
def analyze_image():
    image_url = request.json.get('image_url')
    
    response = openai.ChatCompletion.create(
        model="gpt-4-vision-preview",
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Придумай смешную подпись к этому мему."},
                    {"type": "image_url", "image_url": image_url},
                ],
            }
        ],
        max_tokens=300,
    )
    
    return jsonify({"text": response.choices[0].message.content})

if __name__ == '__main__':
    app.run(port=5000)