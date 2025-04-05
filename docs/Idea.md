# FlashLearn Project Development Summary

## The Idea

FlashLearn is a browser extension designed to streamline the creation of flashcards from online content. The core concept is to eliminate friction between reading content online and creating effective flashcards for later study. Users can highlight text directly in their browser and instantly create flashcards, with AI assistance to suggest potential cards based on the content.

The system will utilize a cloud-based approach for storage, allowing users to access their flashcards across multiple devices. It incorporates spaced repetition techniques for optimized learning and retention.

## Key Considerations

### Technical Development
- **Browser Support**: Initially focusing on Chrome/Chromium browsers due to market share
- **Authentication**: Using OAuth with Google/Apple for seamless user experience
- **Storage**: Cloud-based solution requiring user accounts
- **Integration**: Standalone system with export capability to Anki
- **Tech Stack**: GoLang (backend), JavaScript/TypeScript (extension), React (web interface)
- **Deployment**: Starting with simple Docker containers, with plans to scale using Kubernetes later

### AI Component
- **Model Selection**: Using smaller language models like Gemma-3 4B initially for cost efficiency
- **Privacy**: User data will only be used to improve the product, with opt-out options
- **User Experience**: AI suggestions will be user-triggered rather than automatic
- **Quota System**: Different tiers will offer varying numbers of AI generations per month

### User Experience
- **Card Creation**: Quick highlight-to-card workflow with editing capabilities
- **Card Types**: Support for basic, cloze, and basic-reverse card types initially
- **Review Mechanism**: Separate website for studying flashcards
- **Media Support**: Starting with text-only, with plans to add rich text formatting, images, and audio

## Competitors & Differentiation

### Competitors
- **Anki**: Powerful but requires manual card creation
- **Readwise**: Focuses on note consolidation rather than active recall
- **Remnote**: Standalone app for notes and quizzing
- **Supermemo**: Similar to Anki with high manual effort
- **Quizlet**: Popular flashcard platform
- **Memrise**: Language-focused spaced repetition

### Differentiation
- **Browser Integration**: Creating cards directly while browsing without context switching
- **Reduced Friction**: Eliminating the tedious manual card creation process
- **AI Assistance**: Smart suggestions for card creation based on content
- **Focus on Active Recall**: Unlike Readwise which only shows highlights, FlashLearn focuses on active testing
- **Improved SRS Algorithm**: Plans to enhance the spacing algorithm compared to existing solutions

## Business Model

- **Tiered Subscription**:
  - Basic tier (~$5): No GenAI features, limited storage
  - Standard tier (~$10-15): Full GenAI capabilities with monthly quotas
  - Premium tier: Higher quotas and additional features
- **Add-on Purchases**: Option to buy additional AI generation credits
- **Free Trial**: To encourage user adoption

## Marketing Strategy

- **Development Blog**: Weekly posts on Substack documenting the development journey
- **Development Videos**: YouTube videos showing the creation process and progress
- **Community Building**: Engaging with target audience (students and professionals) throughout development
- **Additional Channels**:
  - Reddit communities focused on productivity and learning (**r/productivity**, **r/studytips**, **r/ADHD**, **r/programming**)
  - Product Hunt launch for MVP
  - Early access waitlist
  - Academic partnerships
  - SEO content around learning techniques

The overall strategy focuses on building in public to gather feedback throughout the development process while simultaneously building an audience of potential users who are invested in the product's success.