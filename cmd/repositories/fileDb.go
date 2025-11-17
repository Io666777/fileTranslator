package repositories

import (
	"fileTranslator/cmd/models"
	"fileTranslator/cmd/storage"
)

func CreateFile(file models.Filestr) (models.Filestr, error) {
    db := storage.GetDB()
    sqlStatement := `INSERT INTO filestrs (namefile, author, cloudkey, langfrom, langto) 
                     VALUES ($1, $2, $3, $4, $5) RETURNING id`
    err := db.QueryRow(sqlStatement, file.NAMEFILE, file.AUTHOR, file.CLOUDKEY, file.LANGFROM, file.LANGTO).Scan(&file.ID)
    if err != nil {
        return file, err
    }
    return file, nil
}