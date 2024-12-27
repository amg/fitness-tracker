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

ARG API_HOST 
ENV API_HOST $API_HOST
RUN echo "API host: $API_HOST"

# Expose the port your app runs on
EXPOSE ${WEB_PORT}

# Make executable
RUN chmod +x ./build_web_env.sh

# Copy some env variables to React script
RUN ./build_web_env.sh

# Define the command to run your app
CMD ["npm", "start"]