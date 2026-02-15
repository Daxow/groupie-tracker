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

        const nameMatch = artist.name && artist.name.toLowerCase().includes(searchValue)

        const membersMatch = artist.members && artist.members.some(function (member) {
                return member.toLowerCase().includes(searchValue)
            })

        const creationMatch = artist.creationDate && artist.creationDate.toString().includes(searchValue)

        const firstAlbumMatch = artist.firstAlbum && artist.firstAlbum.toLowerCase().includes(searchValue)

        const locationMatch = artist.locations && artist.locations.some(function (location) {
                return location.toLowerCase().includes(searchValue)
            })

        return nameMatch || membersMatch || creationMatch || firstAlbumMatch || locationMatch
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