let allArtists = []

fetch("/api/artists")
    .then(response => response.json())
    .then(data => {
        allArtists = data
        console.log(allArtists)
    })