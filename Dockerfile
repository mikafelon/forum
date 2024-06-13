# Utiliser l'image de base golang
FROM golang:1.20

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go.mod et go.sum et télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste des fichiers de l'application
COPY . .

# Lister les fichiers dans le répertoire des templates pour vérification
RUN ls -la /app/templates
RUN ls -la /app/static

# Compiler l'application
RUN go build -o main .

# Exposer le port sur lequel l'application s'exécute
EXPOSE 8080

# Exécuter l'application
CMD ["./main"]
