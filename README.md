<div align="center">
   
   # **Groupie Tracker**
   ### Application web développée en Go permettant d’afficher, rechercher, filtrer et géolocaliser des artistes à partir d’une API externe.
   
</div>

<p align="center">
   <img width="1529" height="895" alt="image" src="https://github.com/user-attachments/assets/2ee7b027-efb1-4839-b045-e1126ffb5835" />
</p>

## **Structure du projet**

```
GROUPIE-TRACKER/
│
├── main.go
├── server.go
├── api.go
├── structs.go
├── go.mod
├── README.md
│
├── templates/
│   ├── index.html
│   └── artist.html
│
└── static/
    ├── style.css
    ├── script.js
    └── artist.js
```

## **Fonctionnalités**
- Affichage des artistes sous forme de cartes
- Page dédiée pour chaque artiste
- Barre de recherche dynamique avec suggestions
- Géolocalisation des concerts

### **Filtres :**
- Date de création
- Premier album
- Nombre de membres
- Lieux de concert

## **Prérequis**
- Go 1.22 ou supérieur installé
- Connexion internet (appel à l’API externe)

## **Installation**

Dans le terminal :

```bash
git clone https://github.com/Daxow/groupie-tracker.git
cd groupie-tracker
```

## **Exécution**

Dans le terminal :

```bash
go run .
```

## **Utilisation**
Une fois le serveur lancé, ouvrir un navigateur et aller à :
```
http://localhost:8080
```
