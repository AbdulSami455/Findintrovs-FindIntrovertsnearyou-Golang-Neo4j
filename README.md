# Find Introvs - Introvert  Matching Platform
Welcome to the Introvert Preferences Matching Platform! This project is designed to create a personalized, introvert-friendly environment where users can connect with like-minded individuals based on their hobbies, preferences, and activities. Using Neo4j as the database and Gin framework for API development, the platform focuses on meaningful connections by matching users based on shared interests.

## Features
## Current Features
### User Registration and Login:

Secure endpoints to register and authenticate users.
### Introvert Preferences Management:
Users can input their preferences for activities like movies, games, music, books, art, and more.
Options to specify likes, dislikes, and interaction preferences.

### Database Insights:
Retrieve database statistics (e.g., node counts) for better data visualization.

### Relationship Management:
Create relationships between nodes to represent connections based on shared preferences.

### Essential Data Addition:
Users can add essential and detailed preference data to refine their matching.
### Introvert-Focused Matching:
Match and assign relationships between users based on attributes and interests stored in the database.

## Technology Stack
Backend Framework: Gin (Go-based web framework)
Database: Neo4j (Graph database for relationship-focused data)
Authentication: Secure endpoints for user management
Networking: RESTful APIs with endpoints for data insertion, retrieval, and relationship creation

## Setup and Installation
Prerequisites
Go: Install Go.
Neo4j: Install and run Neo4j locally or use a cloud instance.
Docker (optional): To containerize the project.
Installation Steps

Clone the Repository:


git clone https://github.com/your-repo-name.git
cd your-repo-name

Install Dependencies:


go mod tidy

Configure Neo4j: Update your Neo4j credentials in the project to connect securely.

Run the Project:


go run main.go
Test Endpoints: Use tools like Postman or cURL to interact with the APIs.

##API Endpoints
### Authentication
POST /api/register: Register a new user.
POST /api/login: Authenticate a user.
POST /api/change-password: Change a user's password.
### Preferences Management
POST /api/nodes: Add essential data nodes.
POST /api/nodes/data: Add introvert-specific preference data.
### Relationship Management
POST /api/relationships: Create relationships between users.
POST /api/match-and-assign-with-attributes: Match and assign relationships based on shared preferences.
Database Insights
GET /api/databases: List all databases.
GET /api/databases/:dbname/count: Get the count of nodes in a specific database.
Future Goals
Real-Time Chat with WebSockets:

Enable instant messaging between matched users using WebSockets.
Create a secure and lightweight chat interface for seamless communication.
Video Calling with WebRTC:

Add real-time video call capabilities to foster deeper connections between users.
Use WebRTC for efficient peer-to-peer communication.
Personalized Bots:

Introduce a personalized chatbot that acts as a companion when users can't find a match.
Use AI to simulate meaningful conversations based on user preferences.
Option to create a bot for users to share with their friends.
Friend Recommendation System:

Suggest potential friends based on shared interests and complementary preferences.
Provide insights into why a recommendation was made.
Enhanced Matching Algorithms:

Refine Neo4j queries to consider more complex relationships and weight-based preferences for improved matches.
