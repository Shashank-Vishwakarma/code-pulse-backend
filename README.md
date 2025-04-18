# Code Pulse - Backend
Code Pulse is a comprehensive online learning platform designed to empower developers in mastering Data Structures, Algorithms, and emerging tech concepts. Built with a focus on interactive learning and skill development, this backend serves as the robust engine powering a seamless coding education experience.

**Key Objectives:**
- Provide an intuitive platform for practicing DSA problems
- Offer in-depth technical blog content for continuous learning
- Create a supportive environment for skill progression
- Enable developers to track and visualize their coding journey

# Features Snapshots
### Landing Page
The landing page introduces the platform and provides options to log in or create a new account for seamless access.

![image](https://github.com/user-attachments/assets/a2c61d42-6f50-43a6-a3a0-6fe7c18cc0f7)

### Online Compiler
Supports Python, JavaScript, Go, C++, and Java for coding directly within the platform.

![image](https://github.com/user-attachments/assets/27c693a6-bbf5-470c-a9cd-7ba2283d08f7)

### Problem Set Page
Displays a comprehensive list of coding problems with a search bar utilizing a debouncing technique for efficient filtering.

![image](https://github.com/user-attachments/assets/1afa105b-57b5-411c-8280-54abab66a29a)

### Problem Description Page
Includes problem title, description, example test cases, hints, topics, a code editor with starter code for Python and JavaScript, and options to run or submit solutions.

![image](https://github.com/user-attachments/assets/11c06b49-1e81-42b6-96e5-5896414395c6)

### Challenges Page
This page has AI-powered generation of challenges and displays unique coding challenges for users to attempt.

![image](https://github.com/user-attachments/assets/9225d5b6-3502-4bd8-a7f2-6ac40eff0e2f)

### Challenge Page
Lists all questions in a selected challenge with a submit button for finalizing solutions.

![image](https://github.com/user-attachments/assets/99012341-c6fd-441b-ae3d-d7f9f6f79694)

### Blogs Page
Showcases all user-published blogs with a search filter that leverages debouncing for quick and precise results.

![image](https://github.com/user-attachments/assets/8bfd98d8-b195-435b-a1c3-4019d78c49b2)

### Blog Page
Displays a banner image, blog title, published date, markdown-formatted content, and an interactive comment section.

![image](https://github.com/user-attachments/assets/7be21fb7-8f2c-4288-ab4c-0a776c23691e)

### Account Profile
Presents user statistics, including total question submissions, questions created, blogs published, challenges created, and challenges attempted.

![image](https://github.com/user-attachments/assets/e5f91f39-3dfd-4365-9186-8b5ebfe33576)

# Features

1. **User Authentication & Management**
   - Secure user registration and login
   - JWT-based authentication
   - User profile management

2. **DSA Question Management**
   - Comprehensive database of Data Structures and Algorithms questions
   - Multiple difficulty levels (Easy, Medium, Hard)
   - Categorization by topics and tags
   - Ability to filter and search questions

3. **Code Submission & Evaluation**
   - Support for multiple programming languages
   - Online code editor with syntax highlighting
   - Automated code compilation and execution
   - Detailed test case results and performance metrics

4. **Progress Tracking**
   - Track solved and attempted questions
   - Performance analytics and progress visualization

5. **Tech Blog Platform**
   - Create, read, and manage tech blog posts
   - Categorization of blogs by technology and topic
   - Commenting and interaction features

6. **AI Quiz Generator**
   - Generate quiz on coding topics such as sql, c++, go, etc using AI

# Tech Stack

## Backend Technologies
1. **Go (Golang)**
   - Primary programming language
   - Known for high performance and concurrency
   - Strong typing and efficient memory management

2. **Gin Web Framework**
   - High-performance HTTP web framework

3. **Database**
   - **MongoDB**
     - NoSQL document database
     - Flexible schema design
     - Horizontal scalability
     - Supports complex queries and indexing
     - Used Atlas Search for text search

4. **Schema Validation**
   - **Go Validator**
     - Standard library for schema validation

5. **Configuration Management**
   - **Viper**
     - Configuration management library
     - Supports multiple file formats (JSON, YAML, ENV)
     - Secure and flexible configuration handling

6. **Authentication**
   - **JWT (JSON Web Tokens)**
     - Secure, stateless authentication mechanism
     - Compact and self-contained
     - Support for claims and token validation

7. **Caching & Performance**
   - **Redis**
     - In-memory data structure store
     - High-performance caching mechanism

8. **Message Queue & Asynchronous Processing**
   - **RabbitMQ**
     - Robust message broker
     - Asynchronous task processing

9. **Email Service**
   - **GoMail**
     - Simple and efficient email sending library

## Additional Tools
- **Postman**: API development and testing
- **VS Code**: Primary development IDE
- **Git**: Version control system
