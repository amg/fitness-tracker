# Use the latest LTS version of Node.js
FROM node:18-alpine
 
# Set the working directory inside the container
WORKDIR /app
 
# Copy package.json and package-lock.json
COPY package*.json ./
 
# Install dependencies
RUN npm install
 
# Copy the rest of your application files
COPY . .

ARG API_BASE_URL
ENV API_BASE_URL=$API_BASE_URL

# Make executable
RUN chmod +x ./run_web.sh

# Run the app
CMD ["./run_web.sh"]