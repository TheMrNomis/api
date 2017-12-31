#HFRlib/API
API pour enregistrer des données de scripts greasemonkey partagées entre plusieurs PC


#Documentation
Dans les exemples qui vont suivre, 2 variables sont à retenir:
* `userID` : une manière unique d'identifier un utilisateur. Tout pourrait être utilisé, mais pour être cohérent entre les différents scripts, merci d'utiliser le hash `bcrypt` du pseudo du posteur.
* `data-name` : le nom de la donnée. Par ex. `blacklist` ou `last_viewed_post`

##Enregistrer des données
```
POST <domain>/script-data/<userID>/<data-name>/
```
Avec, comme `data`, les données à enregistrer

##Lire des données
```
GET <domain>/script-data/<userID>/<data-name>/
```

##Licence
```
GET <domain>/license
```

##Code source
```
GET <domain>/source
```

#Lincense
Affero GNU Public License. Voir LICENSE.md pour les détails.
