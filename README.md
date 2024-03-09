## Project Structure

<pre>
<span style="color: dodgerblue;"><b>./</b></span>
├── .gitignore
├── LICENSE
├── README.md
├── Report.pdf
├── <span style="color: dodgerblue;"><b>assets/</b></span>
├── <span style="color: dodgerblue;"><b>client-mall/</b></span>
├── <span style="color: dodgerblue;"><b>client-whadmin/</b></span>
├── <span style="color: dodgerblue;"><b>server/</b></span>
└── <span style="color: dodgerblue;"><b>database/</b></span>
    ├── <span style="color: dodgerblue;"><b>mysql/</b></span>
    └── <span style="color: dodgerblue;"><b>mongodb/</b></span>
</pre>

1. `.gitignore`: This file tells git which files it should not track / not maintain a version history for.
2. `LICENSE`: MIT License.
3. `README.md`: The current file. Contains the necessary information and instructions to set up and run this project locally, the link to the video demonstration for this project, and contribution information.
4. `Database-Design-and-Implementation-Report.pdf`: The Database Design and Implementation report as required on the [Assessment page on Canvas](https://rmit.instructure.com/courses/121627/assignments/836777).
5. `assets/`: This folder contains the [Project Requirement Specifications](assets/Project_ISYS2099.pdf) PDF file and other miscellaneous medias.
6. `client-mall/`: This NodeJS project directory contains the source code for the Mall frontend client.
7. `client-whadmin/`: This NodeJS project directory contains the source code for the Warehouse Administrator's frontend client.
8. `server/`: This NodeJS project directory contains the source code for the Warehouse Administrator's frontend client.
9. `database/mysql`: This NodeJS project folder contains the setup script (NodeJS), utilities SQL scripts, and sample data for the local MySQL database instance. See **Installation - Database - MySQL** below for more detail.
10. `database/mongodb`: This NodeJS project folder contains the setup script (NodeJS) and sample data for the local MongoDB database instance. See **Installation - Database - MongoDB** below for more detail.


## Technology Stack

- Database: MySQL and MongoDB.
- Backend API Server: NodeJS, Express, and Mongoose.
- Frontend: Vite, React, Bootstrap, and Axios.


## Dependencies

| Dependency                                                                 | Version  |
|:---------------------------------------------------------------------------|:--------:|
| [NodeJS](https://nodejs.org/)                                              | 18.\*.\* |
| [MySQL Community Server](https://dev.mysql.com/downloads/mysql/)           |  8.0.33  |
| [MySQL Shell](https://dev.mysql.com/downloads/shell/)                      |  8.0.33  |
| [MongoDB Community Server](https://www.mongodb.com/try/download/community) |  6.0.\*  |


## Installation

### Database:

#### MySQL:

1. From the project's root directory, navigate to `database/mysql/` directory:
   ```bash
   cd database/mysql
   ```
2. Initialize the MySQL database:

   - Without mock-data:
     ```bash
     npm ci
     npm run setup  # You will be prompted for your MySQL server root username and password
     ```

   - With mock-data (Optional - for testing purpose only):
     ```bash
     npm ci
     npm run setup-with-mock-data  # You will be prompted for your MySQL server root username and password
     ```

   - Alternatively, you can manually set up the MySQL database using MySQL Shell and our provided SQL setup scripts for more granular control:
     - Connect to your local instance of MySQL as the root user:
       ```bash
       mysql -u<root-username> -p
       ```
     - Execute the setup SQL scripts located in the `database/mysql/` directory in the following order:
       - `source reset.sql` (Warning - this will wipe all existing data)
       - `source init/tables.sql`
       - `source init/indexing.sql`
       - `source init/business_rules.sql`
       - `source init/users.sql`
       - `source init/mock_data.sql` (Optional - for testing purpose only)

#### MongoDB:

1. From the project's root directory, navigate to `database/mongdodb/` directory:
   ```bash
   cd database/mongodb
   ```
2. Initialize the MongoDB database:

    - Without mock-data:
      ```bash
      npm ci
      npm run setup
      ```

    - With mock-data (Optional - for testing purpose only):
      ```bash
      npm ci
      npm run setup-with-mock-data
      ```


### Backend API server in `server`:

```bash
cd server
npm clean-install
 ```

### Mall frontend in `client-mall` (for sellers and buyers):

```bash
cd client-mall
npm clean-install
```

### Warehouse dashboard frontend in `client-whadmin` (for warehouse administrators):

```bash
cd client-whadmin
npm clean-install
 ```


## Usage

1. Make sure your local instance of MySQL is running and have been initialized according to the instructions above.
2. Make sure your local instance of MongoDB is running and have been initialized according to the instructions above.
3. Start API server:
   ```bash
   cd server
   npm run start
   ```
   This will start the API server listening at [http://localhost:3000/](http://localhost:3000/)
4. Start Mall frontend:
   ```bash
   cd client-mall
   npm run dev
   ```
   This will start the Mall frontend at [http://localhost:3001/](http://localhost:3001/)
5. Start Warehouse Dashboard frontend:
   ```bash
   cd client-whadmin
   npm run dev
   ```
   This will start the Warehouse Dashboard frontend at [http://localhost:3002/](http://localhost:3002/)


## Video Demonstration:

Available on [YouTube](https://youtu.be/nmXr4WKHXf8).


## Contribution

| SID      | Name                 | Role |
|:---------|:---------------------|:--------------------:|
| s3924826 | Tran Minh Nhat       |   Technical Writer   |
| s3864188 | Phan Thanh Loi       |   Frontend Developer   |
| s3963207 | Do Le Long An        |   Backend Developer   |
| s3877562 | Vo Tuong Minh (Mike) |   Database Admin   |


## Developer Tools

![Git](https://img.shields.io/badge/Git-F05032?style=for-the-badge&logo=git&logoColor=white)
![Vim](https://img.shields.io/badge/Vim-019733?style=for-the-badge&logo=vim&logoColor=white)
![VS Code](https://img.shields.io/badge/VS_Code-0078D4?style=for-the-badge&logo=visual%20studio%20code&logoColor=white)
![WebStorm](https://img.shields.io/badge/WebStorm-000000?style=for-the-badge&logo=webstorm&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)
![MySQL Workbench](https://img.shields.io/badge/MySQL_Workbench-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![MongoDB Compass](https://img.shields.io/badge/MongoDB_Compass-47A248?style=for-the-badge&logo=mongodb&logoColor=white)

# ***TODOs:***
- Separate the FE and BE of the project and deploy it
- Testing the concurrency control when multiple users are buying the same thing
