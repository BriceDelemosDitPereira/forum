# La syntaxe du dockerfile s'écrit comme ceci : 
# COMMANDE/INFO complement

# Origine du langage pour l'image
FROM golang:1.19

# Metadata label
LABEL projet="forum : http.//localhost:8080"
LABEL authors="Thomas DELESTRE &  Brice DELEMOS & Raphael LOVERGNE & Kevin CASTEL & Mickael MARCHAIS"


# Définition de la destination de la copie
WORKDIR /app

# Copie les fichiers spécifiés dans l'image ( "." pour "tous les fichiers")
ADD . /app

# Copie les fichiers spécifiés dans l'image ( "." pour "tous les fichiers" à l'emplacement du fichier dockerfile)
# COPY . /app

# Téléchargement des modules présents dans le go.mod
RUN go mod download

# Build emplacement du main dans l'image + emplacement du main à renseigner
RUN go build -o /forum ./cmd

# Définition du port
EXPOSE 8080

# Création du container depuis l'image [emplacement du main dans l'image]
CMD ["/forum"]
