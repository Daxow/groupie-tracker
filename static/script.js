let allArtists = []

fetch("/api/artists")
    .then(response => response.json())
    .then(data => {
        allArtists = data
        console.log(allArtists)
    })

const searchInput = document.getElementById("searchInput")
const suggestionsContainer = document.getElementById("suggestions")

searchInput.addEventListener("input", function () {

    const searchValue = searchInput.value.toLowerCase()

    suggestionsContainer.innerHTML = ""

    if (searchValue.length === 0) {
        return
    }

    const filteredArtists = allArtists.filter(function (artist) {
        return artist.name.toLowerCase().includes(searchValue)
    })

    filteredArtists.forEach(function (artist) {
        const suggestionElement = document.createElement("div")
        suggestionElement.textContent = artist.name

    suggestionElement.addEventListener("click", function () {
        searchInput.value = artist.name
        suggestionsContainer.innerHTML = ""
    })

    suggestionsContainer.appendChild(suggestionElement)
    })

})