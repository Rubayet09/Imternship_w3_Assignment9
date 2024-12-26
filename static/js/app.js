document.addEventListener('DOMContentLoaded', () => {
    let currentCats = [];
    let currentCatIndex = 0;
    let slideInterval;
    let hasLoadedBreeds = false;

    const tabs = document.querySelectorAll('.tab-btn');
    const sections = document.querySelectorAll('.tab-content');

    // Tab switching logic
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const targetId = tab.getAttribute('data-tab');
            // Clear slideshow when switching away from breeds tab
            if (window.slideInterval && targetId !== 'breeds') {
                clearInterval(window.slideInterval);
                window.slideInterval = null;
            }
            if (targetId === 'breeds' && window.currentBreedImages) {
                startSlideshow();
            }
            // If switching away from breeds tab, clear interval
            else if (window.slideInterval) {
                clearInterval(window.slideInterval);
            }

            tabs.forEach(t => t.classList.remove('active'));
            sections.forEach(s => s.classList.remove('active'));

            tab.classList.add('active');
            document.getElementById(targetId).classList.add('active');

            if (targetId === 'voting' && currentCats.length === 0) {
                loadCats();
            } else if (targetId === 'breeds') {
                loadBreeds();
            } else if (targetId === 'favs') {
                loadFavorites();
            }
        });
    });

    // Voting section
    async function loadCats() {
        try {
            const response = await fetch('/api/cats');
            const data = await response.json();
            if (data.status === 'success') {
                currentCats = data.data;
                currentCatIndex = 0;
                displayCurrentCat();
            }
        } catch (error) {
            console.error('Error loading cats:', error);
        }
    }

    function displayCurrentCat() {
        if (currentCats.length === 0) return;

        const catImg = document.querySelector('.cat-image img');
        if (catImg) {
            catImg.src = currentCats[currentCatIndex].url;
            console.log('Displaying image:', currentCats[currentCatIndex].url); // Debug log

            // Add error handling for image load
            catImg.onerror = () => {
                console.error('Failed to load image:', currentCats[currentCatIndex].url);
                catImg.src = 'path/to/fallback-image.jpg'; // Add a fallback image
            };
        }
    }


    // In app.js - Update the voteCat function
    async function voteCat(vote) {
        try {
            const cat = currentCats[currentCatIndex];
            if (!cat || !cat.id) {
                console.error('Invalid cat data');
                return;
            }

            const response = await fetch('/api/vote', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    image_id: cat.id,
                    image_url: cat.url,
                    vote: vote
                })
            });

            const data = await response.json();
            if (data.status !== 'success') {
                throw new Error(data.message);
            }

            currentCatIndex++;
            if (currentCatIndex >= currentCats.length) {
                await loadCats();
            } else {
                displayCurrentCat();
            }
        } catch (error) {
            console.error('Error voting:', error);
        }
    }

    async function loadBreedsWithDefault() {
        try {
            const response = await fetch('/api/breeds');
            const data = await response.json();

            if (data.status === 'success' && Array.isArray(data.data)) {
                const select = document.getElementById('breed-select');
                select.innerHTML = '<option value="">Select a breed</option>';

                data.data.sort((a, b) => (a.name || a.Name).localeCompare(b.name || b.Name));

                let abyssinianId = null;

                data.data.forEach(breed => {
                    const option = document.createElement('option');
                    option.value = breed.id || breed.ID;
                    option.textContent = breed.name || breed.Name;
                    select.appendChild(option);

                    // Store Abyssinian ID if found
                    if ((breed.name || breed.Name).toLowerCase() === 'abyssinian') {
                        abyssinianId = breed.id || breed.ID;
                    }
                });

                // Add change event listener
                select.removeEventListener('change', onBreedSelect);
                select.addEventListener('change', onBreedSelect);

                // Set Abyssinian as default if found
                if (abyssinianId) {
                    select.value = abyssinianId;
                    loadBreedDetails(abyssinianId);
                }

                hasLoadedBreeds = true;
            }
        } catch (error) {
            console.error('Error loading breeds:', error);
        }
    }



    async function loadBreeds() {
        try {
            const response = await fetch('/api/breeds');
            const data = await response.json();

            if (data.status === 'success' && Array.isArray(data.data)) {
                const select = document.getElementById('breed-select');
                select.innerHTML = '<option value="">Select a breed</option>';

                data.data.sort((a, b) => (a.name || a.Name).localeCompare(b.name || b.Name));

                data.data.forEach(breed => {
                    const option = document.createElement('option');
                    option.value = breed.id || breed.ID;
                    option.textContent = breed.name || breed.Name;
                    select.appendChild(option);
                });

                select.removeEventListener('change', onBreedSelect);
                select.addEventListener('change', onBreedSelect);
            }
        } catch (error) {
            console.error('Error loading breeds:', error);
        }
    }

    async function onBreedSelect(e) {
        const breedId = e.target.value;
        if (!breedId) {
            document.querySelector('.breed-details').style.display = 'none';
            return;
        }

        await loadBreedDetails(breedId);
    }

    async function loadBreedDetails(breedId) {
        const loadingEl = document.getElementById('breed-loading');
        const breedDetails = document.querySelector('.breed-details');

        loadingEl.classList.remove('hidden');
        breedDetails.style.display = 'none';

        if (window.slideInterval) {
            clearInterval(window.slideInterval);
        }

        try {
            const response = await fetch(`/api/breed?id=${breedId}`);
            const data = await response.json();

            if (data.status === 'success' && data.data) {
                const breed = data.data;
                console.log('Breed data:', breed); // Log breed data
                console.log('Breed images:', breed.images); // Log images specifically

                window.currentBreedImages = breed.images || [];
                window.currentImageIndex = 0;

                updateBreedDisplay(breed); // Pass the breed object
                loadingEl.classList.add('hidden');
                breedDetails.style.display = 'block';
            }
        } catch (error) {
            console.error('Error loading breed details:', error);
            loadingEl.textContent = 'Error loading breed details';
        }
    }



    function updateBreedDisplay(breed) {
        const breedTitle = document.querySelector('.breed-title');
        const breedOrigin = document.querySelector('.breed-origin');
        const breedDesc = document.querySelector('.breed-description');
        const breedWiki = document.querySelector('.wiki-link');
        const breedImage = document.querySelector('.breed-image'); // Single <img> element
        const dotsContainer = document.querySelector('.slide-dots');
        const images = breed.images || []; // Fetch images array

        console.log('Breed data received:', breed);
        console.log('Breed images:', images); // Log images for debugging

        if (!breedTitle || !breedImage || !dotsContainer) {
            console.error('Required elements are missing in the DOM.');
            return;
        }

        // Update text content
        breedTitle.textContent = breed.name || 'Unknown Breed';
        breedOrigin.textContent = breed.origin ? `(${breed.origin})` : '';
        breedDesc.textContent = breed.description || 'Description not available.';

        // Update Wikipedia link
        if (breed.wikipedia_url) {
            breedWiki.href = breed.wikipedia_url;
            breedWiki.style.display = 'inline-block';
            breedWiki.textContent = 'Wikipedia';
        } else {
            breedWiki.style.display = 'none';
        }

        // Update the breed image
        if (images.length > 0) {
            console.log('Setting image src to:', images[0]); 
            
            breedImage.style.display = 'block';// Debugging log
            breedImage.src = images[0];
            breedImage.alt = breed.name || 'Breed image';
            breedImage.onerror = () => {
                console.error('Failed to load image:', images[0]);
                breedImage.src = 'https://via.placeholder.com/300'; // Replace with a valid fallback image
            };

            // Generate dots for image navigation
            dotsContainer.innerHTML = images.map((_, idx) => `
                    <div class="dot ${idx === 0 ? 'active' : ''}" data-index="${idx}"></div>
                `).join('');

            // Add click events for dots
            dotsContainer.querySelectorAll('.dot').forEach((dot, idx) => {
                dot.addEventListener('click', () => {
                    changeImage(idx, images); // Handle dot click
                });
            });

            // Start the slideshow
            window.currentBreedImages = images;
            window.currentImageIndex = 0;
            startSlideshow();

            const noImageMessage = document.querySelector('.no-image-message');
            if (noImageMessage) {
                noImageMessage.remove();
            }


        } else {
            // No images available, show a message and hide the image container
            console.error('No images available for this breed.');
            breedImage.style.display = 'none'; // Hide the image container
            breedImage.src = ''; // Clear the image
            breedImage.alt = 'No image available for this breed'; // Set alt text
    
            // Display "No image" message if not already shown
            let noImageMessage = document.querySelector('.no-image-message');
            if (!noImageMessage) {
                noImageMessage = document.createElement('p');
                noImageMessage.className = 'no-image-message';
                noImageMessage.textContent = "There's no image for this breed.";
                noImageMessage.style.color = 'gray';
                breedImage.parentElement.appendChild(noImageMessage); // Add the message below the image container
            }

        // Clear the dots container as there are no images to navigate
            dotsContainer.innerHTML = '';
        }
    }





    function startSlideshow() {
        // Clear any existing interval
        if (window.slideInterval) {
            clearInterval(window.slideInterval);
        }

        const images = window.currentBreedImages;
        if (images && images.length > 1) {
            window.slideInterval = setInterval(() => {
                if (document.getElementById('breeds').classList.contains('active')) {
                    const nextIndex = (window.currentImageIndex + 1) % images.length;
                    changeImage(nextIndex, images);
                }
            }, 3000);
        }
    }




    function changeImage(index, images) {
        if (!images || !images[index]) {
            console.error('Invalid image index:', index);
            return;
        }

        window.currentImageIndex = index;

        // Update image
        const breedImage = document.querySelector('.breed-image');
        console.log('Changing to image:', images[index]); // Debugging log
        breedImage.src = images[index];
        breedImage.alt = `Breed image ${index + 1}`;

        // Update active dot
        document.querySelectorAll('.dot').forEach((dot, idx) => {
            dot.classList.toggle('active', idx === index);
        });

        // Reset slideshow timer
        startSlideshow(images);
    }





    // Make changeImage available globally
    window.changeImage = changeImage;

    const breedsTab = document.querySelector('[data-tab="breeds"]');
    breedsTab.addEventListener('click', () => {
        if (!hasLoadedBreeds) {
            loadBreedsWithDefault();
        } else {
            loadBreeds();
        }
    });

    // Initialize breeds when the breeds tab is clicked
    // const breedsTab = document.querySelector('[data-tab="breeds"]');
    // breedsTab.addEventListener('click', loadBreeds);



    // Favorites section
    // Update the loadFavorites function in your app.js

    async function loadFavorites() {
        try {
            const response = await fetch('/api/favorites');
            const data = await response.json();

            if (data.status === 'success') {
                const grid = document.querySelector('.favorites-grid');
                grid.innerHTML = ''; // Clear existing items

                if (data.data.length === 0) {
                    // Show message if no favorites
                    const message = document.createElement('p');
                    message.textContent = 'No favorite cats yet! Click the ❤️ button while voting to add some.';
                    message.style.gridColumn = '1 / -1';
                    message.style.textAlign = 'center';
                    message.style.color = '#666';
                    grid.appendChild(message);
                    return;
                }

                data.data.forEach(favorite => {
                    const itemDiv = document.createElement('div');
                    itemDiv.className = 'favorite-item';

                    const img = document.createElement('img');
                    img.src = favorite.url;
                    img.alt = 'Favorite cat';

                    const removeBtn = document.createElement('button');
                    removeBtn.className = 'remove-favorite';
                    removeBtn.innerHTML = '×';
                    removeBtn.title = 'Remove from favorites';

                    // In app.js - Update the remove favorite event listener
                    removeBtn.addEventListener('click', async (e) => {
                        e.preventDefault();
                        e.stopPropagation();

                        try {
                            // Log the ID being sent for deletion
                            console.log('Attempting to delete favorite with ID:', favorite.id);

                            const response = await fetch(`/api/favorites/${encodeURIComponent(favorite.id)}`, {
                                method: 'DELETE'
                            });

                            const result = await response.json();
                            if (result.status === 'success') {
                                itemDiv.style.opacity = '0';
                                setTimeout(() => {
                                    itemDiv.remove();
                                    if (grid.children.length === 0) {
                                        loadFavorites();
                                    }
                                }, 300);
                            } else {
                                console.error('Failed to remove favorite:', result.message);
                            }
                        } catch (error) {
                            console.error('Error removing favorite:', error);
                        }
                    });

                    itemDiv.appendChild(img);
                    itemDiv.appendChild(removeBtn);
                    grid.appendChild(itemDiv);

                    img.onerror = () => {
                        itemDiv.remove();
                        console.error('Failed to load favorite image:', favorite.url);
                    };
                });
            }
        } catch (error) {
            console.error('Error loading favorites:', error);
        }
    }

    // Setup voting buttons
    document.querySelector('.vote-btn.like').addEventListener('click', () => voteCat('like'));
    document.querySelector('.vote-btn.dislike').addEventListener('click', () => voteCat('dislike'));
    document.querySelector('.vote-btn.favorite').addEventListener('click', () => voteCat('love'));

    // Initial load
    loadCats();
});