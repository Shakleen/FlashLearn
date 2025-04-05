# Software Product Requirements Document (PRD)

## 1. Introduction

### 1.1. Purpose

This document outlines the requirements for the product “FlashLearn”.

### 1.2. Scope

FlashLearn (FL) is a subscription-based online software. It’s a web extension that allows user to create flash cards while consuming the information intuitively. Users can manually create them or utilize the power of Generative AI to automatically create helpful flash cards. Afterwards, the cards can be studied on a web interface at the dedicated website.

### 1.3. Definitions, Acronyms, and Abbreviations

* **FL (FlashLearn)**: The name of the product.

* **GenAI (Generative AI)**: Large Language Models and Small Language Models used to generate text or image or audio when prompted based on context.

* **SRS (Spaced Repetition System)**: A research based effective study mechanism where hard to learn cards are shown more frequently than easy to learn ones. Used by popular apps like Anki, Quizlet, etc.

* **CICD (Continuous Integration and Continuous Development)**: A development strategy incorporating continuously testing developed code to easily catch bugs as they happen.

### 1.4. References

* Use cases covered in use case diagrams under [Appendix](#6-appendix-optional)

## 2. Overall Description

### 2.1. Product Perspective:

Standalone product with future integrations to other products yet to be announced.

### 2.2. Product Functions:

* Minimize friction of flash card creation through browser extension and integration of GenAI.

* Study created flash cards using SRS in online website accessed through any web enabled device.

* Export created flash cards to other applications like Anki.

### 2.3. User Classes and Characteristics:

* **Students**: Anyone who is currently pursuing formal or informal education which requires consumption and recall of information. This includes but is not limited to high school students, college students, etc.

* **Professionals**: Anyone who is working and needs to retain information about their work, like doctors staying up-to-date on latest advancements in the medical field, Software Engineers studying about the latest technology, etc.

### 2.4. Operating Environment:

* **Hardware requirements**: From a user standpoint, only an internet enabled device. From a developer standpoint, databases, servers, and compute for AI models.

* **Software requirements**: Users need only a browser to access the product. Primarily chromium-based browsers will be supported.

### 2.5. Design and Implementation Constraints:

* **Compute Constraint**: MVP will be using a free tier of cloud providers like GCP or AWS or Heroku. Large AI models aren’t feasible.

* **Performance Considerations**: Develop backend using GoLang such that it’s performant.

### 2.6. User Documentation:

* **Web Extension**: Text based documentation should be enough. So a simple website.

* **Web page**: Video tutorials.

## 3. Specific Requirements

### 3.1. Functional Requirements:

* **Feature/Use Case 1: Create flash cards in browser**

    * User must be able to create flash cards while reading articles online. When a user highlights a text block using their cursor, an option pops-up to convert that into a flash card either manually or using GenAI.

    * Logic: TBD

* **Feature/Use Case 2: Study flash cards in website**

    * User can study created flash cards when they visit the flash learn website. Each study session is focused on a specific deck. Flash cards should be scheduled using SRS initially and then modified SRS powered by a lightweight recommendation system.

    * Logic: TBD

* **Feature/Use Case 3: Manage Decks**

    * User must be able to create and maintain decks. A deck is a collection of cards grouped together by the user. For example, a user can have a deck for “data structures and algorithms” which will contain flash cards for DSA, while another deck for “system design” for system design related cards.

    * Logic: TBD

* **Feature/Use Case 4: Tag flash cards**

    * A user can assign one or more tags to each flash cards. Tags correspond to specific topic while deck corresponds to more general subject. For example, a flash card can have tags “data\_structure::tree::BFS” while being in the deck “Data Structure”.

    * Logic: TBD

* **Feature/Use Case 5: Study statistics**

    * User must be able to view their study statistics. This includes seeing basic stats like how long they study each day, when they study, and what are the different levels of their cards. Later on, more useful analytics should be displayed. For example, what sort of information is the user good at retaining, what helps the user retain information (mnemonics or images or audio cues etc.)

    * Logic: TBD

* **Feature/Use Case 6: Export data**

    * Users must have the ability to export created flash cards into a common format that can be ingested by other popular apps like Anki, Quizlet, etc.

    * Logic: TBD

* **Feature/Use Case 7: Custom study**

    * Users should be able to schedule custom study sessions outside their regular studies. This could include studying past x days worth of new cards, studying forgotten cards, studying cards belonging to a specific tag, etc.

    * Logic: TBD

    * More features: TBD

### 3.2. Non-Functional Requirements

* **3.2.1. Performance Requirements:**

    * Response time: TBD

    * Throughput: TBD

    * Resource utilization: TBD

* **3.2.2. Security Requirements:**

    * Authentication: TBD

    * Authorization: TBD

    * Data encryption: TBD

    * Vulnerability management: TBD

* **3.2.3. Usability Requirements:**

    * Ease of use: Must lessen the friction of card creation as much as possible. Prefer simplicity over complexity.

    * Learnability: Must have a low learning curve.

    * Accessibility: Must be accessible by anyone who has a browser (Initially Chrome browser and then all chromium based browsers).

* **3.2.4. Reliability Requirements:**

    * Availability: Only available while online. Offline functionalities could include caching created flash cards for later syncing.

    * Fault tolerance: TBD

    * Recovery: TBD

* **3.2.5. Maintainability Requirements:**

    * Code quality: Simple, testable, and maintainable codebase.

    * Modularity: Prefer simplicity than over-abstraction of concepts.

    * Documentation: Each and every function must have docstring.

    * CICD: Must incorporate CICD pipeline. Block merges to main if all tests don’t pass.

* **3.2.6. Portability Requirements:**

    * Operating system compatibility: Windows, macOS, and Linux.

    * Browser compatibility: Chrome initially. Then any chromium-based browser. Then others.

* **3.2.7. Scalability Requirements:**

    * TBD

### 3.3. Interface Requirements

* **3.3.1. User Interfaces:**

    * \[Description of the user interface, including layout, navigation, and interaction.\]

* **3.3.2. Hardware Interfaces:**

    * \[Description of any hardware interfaces, such as sensors or peripherals.\]

* **3.3.3. Software Interfaces:**

    * \[Description of any software interfaces, such as APIs or databases.\]

* **3.3.4. Communication Interfaces:**

    * \[Description of any communication interfaces, such as network protocols.\]

### 3.4. Data Requirements

* User data will be stored primarily in structured DBMS format.

* Design of database: TBD

## 4. Acceptance Criteria

* \[Define the criteria that will be used to determine whether the software product meets the requirements. These should be specific, measurable, achievable, relevant, and time-bound (SMART).\]

* \[Example : "The system shall load the user dashboard in less than 2 seconds, 95% of the time."\]

## 5. Future Considerations (Optional)

* GenAI to judge correctness of answer and provide feedback as to what was lacking.

* OCR to extract text from physic documents and books.

* Support for video and audio learning media.

## 6. Appendix (Optional)

* Use case diagrams
    1. [Browser extension](diagrams/usecase_browser_extension.png)
    1. [Website](diagrams/usecase_website.png)
