let allArtists = []

fetch("/api/artists")
    .then(response => response.json())
    .then(data => {
        allArtists = data
    })

const searchInput = document.getElementById("searchInput")
const suggestionsContainer = document.getElementById("suggestions")

searchInput.addEventListener("input", function () {

    const searchValue = searchInput.value.toLowerCase().trim()

    suggestionsContainer.innerHTML = ""
    suggestionsContainer.style.display = "none"

    if (searchValue.length === 0) {
        return
    }

    const suggestions = new Set()
    
    allArtists.forEach(function (artist) {

        if (artist.name && artist.name.toLowerCase().includes(searchValue)) {
            suggestions.add(artist.name + " - Artiste")
        }

        if (artist.members) {
            artist.members.forEach(function (member) {
                if (member.toLowerCase().includes(searchValue)) {
                    suggestions.add(member + " - Membre")
                }
            })
        }

        if (artist.creationDate && artist.creationDate.toString().includes(searchValue)) {
            suggestions.add(artist.creationDate + " - Date de création")
        }

        if (artist.firstAlbum && artist.firstAlbum.toLowerCase().includes(searchValue)) {
            suggestions.add(artist.firstAlbum + " - Premier album")
        }

        if (artist.locations) {
            artist.locations.forEach(function (location) {
                if (location.toLowerCase().includes(searchValue)) {
                    suggestions.add(location + " - Lieu de concert")
                }
            })
        }
    })

    suggestions.forEach(function (text) {

        const suggestionElement = document.createElement("div")
        suggestionElement.textContent = text

        suggestionElement.addEventListener("click", function () {
            searchInput.value = text.split(" - ")[0]
            suggestionsContainer.innerHTML = ""
            searchInput.form.submit()
        })

        suggestionsContainer.appendChild(suggestionElement)
        if (suggestions.size > 0) {
            suggestionsContainer.style.display = "block"
        }
    })

})