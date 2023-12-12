# comeback-ai-bot
Bonjour

# Control flow diagram
https://drive.google.com/file/d/1POUD9P0RiGmMy52zmd_g2dLO4G7GJYUQ/view?usp=sharing

### User Input Query
- User input in slack channel with the query decoded via request param.
- Send the raw statement to wit ai to construct wolfram understandable query by NLP.
- Parse the response and send the query string to wolfram and respond with the answer to slack thread.

### Age Calculator
- User input with the year with command "my yob is 1990".
- Convert to integer then calculate age and response reply back.

### File Upload
- Upload a specified file in the same directory like sample.csv
- Use FileUploadParameters and UploadFile method from official slack go package.
