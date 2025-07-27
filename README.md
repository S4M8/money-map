# 🗺️ Money Map: Your Personal Finance Adventure! 🗺️

Welcome to Money Map, the fun and easy way to navigate your financial world! 🧭 Track your income, manage your expenses, and chart a course to financial freedom. Let's make budgeting an adventure! 🚀

## ✨ Key Features

-   **📊 Interactive Dashboard:** Get a clear, at-a-glance view of your financial landscape.
-   **💸 Track Income & Expenses:** Easily add and manage your capital (income) and your "Core" & "Choice" expenses.
-   **💯 Money Map Score:** Get a simple, intuitive score (from "Poor" to "Great" 👍) that tells you how well you're sticking to the 50/30/20 rule.
-   **💰 Fund Allocation:** Automatically distribute your savings into different funds like an Emergency Fund, Education Fund, and Investments.
-   **🐳 Dockerized:** The entire application is containerized, making setup a breeze!

## 🚀 Getting Started

Ready to start your adventure? You'll just need [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/money-map.git
    cd money-map
    ```

2.  **Create your `.env` file:**
    Copy the example environment file. The defaults are already set up for you!
    ```sh
    cp .env.example .env
    ```

3.  **Launch the application! 🚀**
    Use the included `Makefile` to build and run the application with a single command:
    ```sh
    make up
    ```
    This will start the frontend, backend, and database containers.

4.  **Explore Your Money Map!**
    Once the containers are running, open your browser and navigate to [http://localhost:8080](http://localhost:8080).

## 🛠️ Useful `make` Commands

-   `make up`: Starts the application.
-   `make down`: Stops the application.
-   `make build`: Builds the application images.
-   `make rebuild`: Rebuilds the application from scratch (great for when you make changes!).
-   `make logs`: View the application logs.
-   `make clean`: Stops the application and removes all data.

Happy budgeting, adventurer! 🎉
