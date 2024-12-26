
# Internship_w3_Assignment9


## Cat Voting App- Web Application_Internship Project A9

## Table of Contents 
1. [Overview](#overview) 
2. [Key Features](#key-features) 
3.  [Technologies Used](#technologies-used) 
4. [Development Setup](#development-setup) 
5. [Project Structure](#project-structure) 
6. [Author](#author)


## Overview
The Cat Voting App is a web application built with Beego framework that integrates with The Cat API to allow users to view and vote on cat images accordingly, view images and description with wiki-links of different breeds of cats, also allows user to see their chosen favorite cat images with a remove option. The application implements a concurrent processing model using Go channels for API requests and features a responsive frontend interface.



## Key Features 

- **Image Voting System**: Browse random cat images, Vote with like/dislike/love options, Add favorites with heart button, Automatic image cycling after voting.
- **Breed Explorer**: Fetch all of cat breeds, Detailed breed information including: Origin, Description, Wikipedia link. Also, Default selection of 'Abyssinian' breed, Image slideshow for each breed, Auto-cycling breed images with manual navigation.
 - **Favorites Management**: Save favorite cat images, View all favorites in a grid layout, Remove favorites with confirmation.




## Technologies Used

- **Framework**: Golang [Beego Framework v2.x]
- **Front-End**: JavaScripts[Vanilla JS], HTML5, CSS3
- **Testing**: Unit testing with Go's testing package
- **Containerization**: Docker (optional for PostgreSQL setup)


## Development Setup

### Prerequisites
 Ensure the following are installed on your system:

- go version >=  1.16
- bee v2.3.0

### Step: 1. Clone the repository:
 ```bash 
git clone https://github.com/Rubayet09/Internship_w3_Assignment9.git 
cd Internship_w3_Assignment9
 ```

### Step: 2. Install dependencies:
 ```bash 
go mod tidy
 ```

### Step: 3. Run application:
 ```bash 
bee run
 ```

## How to Use
1. Navigate to http://localhost:8080
2. Voting page should appear where you should be able to POST you votes- like/dislike/love
3. Go to Breeds tab and default breed 'Abyssinian' would appear, you can select your desired breed from the drop down bar from where all the breeds information gets fetched. The images will slide automatically.
4. Next, go to the favs tab where you can see all the images that you voted love to. You can also remove any image from your favorite list of images.

### Steps: 6. Running Tests
```bash 
go test ./models -v -coverprofile=coverage.out

go test ./controllers -v -coverprofile=coverage.out 

go test ./routers -v -coverprofile=coverage.out
 ```
 
NB: After running the command-
models--> coverage: 84.2% of statements
controllers--> coverage: 84.1% of statements
routers--> coverage: 100.0% of statements


### API Integration
---
- You can retrieve all the breeds of the catapi from here- http://localhost:8080/api/breeds
- You can retrieve any specific cat breed for example: 'Bombay', the id of this breed is 'bomb', from here- http://localhost:8080/api/breed?id=bomb
- You can retrieve all your favorite images from here- http://localhost:8080/api/favorites




---
## Project Structure
```
CATVOTINGAPP/
├── conf/
│   └── app.conf          
├── controllers/
│   ├── cat_controller.go 
│   └── cat_controller_test.go
├── models/
│   ├── cat_model.go    
│   └── cat_model_test.go
├── routers/
│   ├── router_test.go
│   └── router.go     
├── static/
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js       
├── views/
│   └── index.tpl            
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Author

Rubayet Shareen
SWE Intern, W3 Engineers
Dhaka, Bangladesh
