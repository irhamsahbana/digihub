# Digihub Backend

## Project Overview

### Problem Statement

In the automotive industry, dealerships often face operational inefficiencies due to disconnected systems and manual processes. Key challenges include:

* **Inefficient Vehicle Inspections:** Manual walk-around checks are slow, inconsistent, and difficult to track, leading to delays and potential oversights.
* **Fragmented Data:** Information on clients, vehicles, promotions, and sales activities is often scattered across different platforms, hindering a unified view of the business.
* **Limited Performance Insight:** Without a central data hub, it's difficult for management to monitor key performance indicators and make informed, strategic decisions.

### The Proposed Solution

**Digihub** is a comprehensive backend system engineered to digitalize and streamline automotive dealership operations. It acts as a central nervous system, providing a suite of robust APIs to unify data, automate workflows, and enhance communication.

By connecting previously disparate modules, Digihub empowers dealerships to:

* Improve the speed and quality of vehicle inspections.
* Maintain a single source of truth for all operational data.
* Securely manage user access with a flexible role-based system.
* Enable activity tracking.

### Key Features

* **Walk-Around Check (WAC) Management:** A complete digital workflow for conducting vehicle inspections, including condition tracking, follow-up logs, and activity monitoring.
* **Master Data Management:** Centralized control over core business entities like Branches, Vehicle Types, and Promotions.
* **User & Client Management:** A secure system for managing employee and customer data.
* **Role-Based Access Control (RBAC):** Ensures users can only access features and data appropriate for their roles.
* **Secure Authentication:** Features a robust authentication system using JWT, including email verification and password reset capabilities.
* **Local File Storage:** Store and manage files securely within the application’s local storage.
* **Application Monitoring:** Provides comprehensive application logging for tracking system health and user activities.

### Tech Stack

* **Language & Framework:** **Go** with the **Fiber** framework
* **Database:** **PostgreSQL** with **sqlx**
* **Authentication:** **JSON Web Token (JWT)**
* **Logging:** **Zerolog**
* **Excel Processing:** **Excelize** (`github.com/xuri/excelize/v2`)
* **Task Runner:** **Go-Task**

## Developer Guide

### General Rules

* We use conventional commits to deal with git commits: <https://www.conventionalcommits.org>
  * Use `feat: commit message` to do git commit related to feature.
  * Use `refactor: commit message` to do git commit related to code refactorings.
  * Use `fix: commit message` to do git commit related to bugfix.
  * Use `test: commit message` to do git commit related to test files.
  * Use `docs: commit message` to do git commit related to documentations (including README.md files).
  * Use `style: commit message` to do git commit related to code style.

* Use git-chglog <https://github.com/git-chglog/git-chglog> to generate changelog (CHANGELOG.md) before merging to release branch.

### Branching Strategy

* Keep your branch strategy simple. Build your strategy from these three concepts:
  * Use feature branches for all new features and bug fixes.
  * Merge feature branches into the main branch using pull requests.
  * Keep a high quality, up-to-date main branch.

#### Use feature branches for your work

Develop your features and fix bugs in feature branches based off your main branch. These branches are also known as
topic branches. Feature branches isolate work in progress from the completed work in the main branch. Git branches are
inexpensive to create and maintain. Even small fixes and changes should have their own feature branch.

<!-- <p align="left"><img src="./featurebranching.png" width="360"></p> -->

#### Name your feature branches by convention

* Use a consistent naming convention for your feature branches to identify the work done in the branch. You can also
  include other information in the branch name, such as who created the branch.

* Some suggestions for naming your feature branches:
  * users/username/description
  * users/username/workitem
  * bugfix/description
  * feature/feature-name
  * feature/feature-area/feature-name
  * hotfix/description

#### Use release branches

* Create a release branch from the main branch when you get close to your release or other milestone, such as the end of
  a sprint. Give this branch a clear name associating it with the release, for example release/20.
* Create branches to fix bugs from the release branch and merge them back into the release branch in a pull request

<!-- <p align="left"><img src="./releasebranching_release.png" width="360"></p> -->
