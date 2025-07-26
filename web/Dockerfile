# Use an official Node.js runtime as a parent image
FROM node:22-alpine

# Set the working directory inside the container
WORKDIR /app


# Copy package.json and package-lock.json (if available)
COPY package*.json ./
COPY .npmrc ./
# Install dependencies
RUN npm ci -f

RUN npm i -g serve
# Copy the frontend source code
COPY . .

# Build the React application with Vite
RUN npm run build


# Expose port 5173 for the Vite preview server
EXPOSE 5173

# Serve the built application using Vite's preview server
CMD ["serve", "-s", "dist"]
