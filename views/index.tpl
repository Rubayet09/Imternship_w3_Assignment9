<!DOCTYPE html>
<html>
<head>
    <title>Cat Voting App</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <script src="/static/js/app.js" defer></script>
</head>
<body>
    <div class="app-container">
        <nav class="tabs">
            <button class="tab-btn active" data-tab="voting">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path d="M9 17.586V3H7v14.586l-2.293-2.293-1.414 1.414L8 21.414l4.707-4.707-1.414-1.414L9 17.586zM20.707 7.293 16 2.586l-4.707 4.707 1.414 1.414L15 6.414V21h2V6.414l2.293 2.293 1.414-1.414z"/>
                </svg>
                <span>Voting</span>
            </button>
            <button class="tab-btn" data-tab="breeds">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                </svg>
                <span>Breeds</span>
            </button>
            <button class="tab-btn" data-tab="favs">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" fill="currentColor"/>
                </svg>
                <span>Favs</span>
            </button>
        </nav>

        <main class="content">
            <!-- Voting section remains unchanged -->
            <section id="voting" class="tab-content active">
                <div class="voting-container">
                    <div class="cat-image">
                        <img src="" alt="Cat">
                    </div>
                    <div class="voting-buttons">
                        <button class="vote-btn dislike">👎</button>
                        <button class="vote-btn favorite">❤️</button>
                        <button class="vote-btn like">👍</button>
                    </div>
                </div>
            </section>

            <!-- Modified Breeds section -->
            <!-- Replace the breeds section in index.tpl with this -->
            <section id="breeds" class="tab-content">
                <div class="breeds-container">
                    <div class="breed-select-wrapper">
                        <select id="breed-select">
                            <option value="">Select a breed</option>
                        </select>
                    </div>
                    
                    <!-- Add loading indicator -->
                    <div id="breed-loading" class="hidden">Loading...</div>
                    
                    <!-- Breed display area -->
                    <div class="breed-details" style="display: none;">
                        <div class="breed-image-container">
                            <img src="" alt="" class="breed-image">
                            <div class="slide-dots"></div>
                        </div>
                        
                        <div class="breed-info">
                            <h2 class="breed-title"></h2>
                            <span class="breed-origin"></span>
                            <p class="breed-description"></p>
                            <a href="" target="_blank" class="wiki-link">WIKIPEDIA</a>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Modified Favorites section -->
            <section id="favs" class="tab-content">
                <div class="favorites-container">
                    <h2>Your Favorite Cats</h2>
                    <div class="favorites-grid">
                        <!-- Template for favorite items -->
                        <template id="favorite-template">
                            <div class="favorite-item">
                                <img src="" alt="Favorite cat">
                                <button class="remove-favorite" title="Remove from favorites">×</button>
                            </div>
                        </template>
                    </div>
                </div>
            </section>
        </main>
    </div>
</body>
</html>