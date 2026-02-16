document.addEventListener("DOMContentLoaded", function () {

    const urlParams = new URLSearchParams(window.location.search)
    const artistId = urlParams.get("id")

    const map = L.map("map").setView([20, 0], 2)

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: "© OpenStreetMap"
    }).addTo(map)

    fetch("/api/artist-locations?id=" + artistId)
        .then(response => response.json())
        .then(data => {

            const markers = []

            data.forEach(coord => {
                const marker = L.marker([coord.lat, coord.lon]).addTo(map)
                markers.push(marker)
            })

            if (markers.length > 0) {
                const group = new L.featureGroup(markers)
                map.fitBounds(group.getBounds())
            }

        })

})
