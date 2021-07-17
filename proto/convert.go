package proto

import (
	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/br3w0r/gamelist-backend/helpers"
)

func (g *GameProperties) ConvertToEntity() entity.GameProperties {
	nPlatforms := len(g.Platforms)
	nGenres := len(g.Genres)

	platforms := make([]entity.Platform, nPlatforms)
	genres := make([]entity.Genre, nGenres)

	maxIter := helpers.MaxInt(nPlatforms, nGenres)
	for i := 0; i < maxIter; i++ {
		if i < nPlatforms {
			platforms[i] = entity.Platform{Name: g.Platforms[i].Name}
		}
		if i < nGenres {
			genres[i] = entity.Genre{Name: g.Genres[i].Name}
		}
	}

	return entity.GameProperties{
		Name:         g.Name,
		Platforms:    platforms,
		YearReleased: uint16(g.YearReleased),
		ImageURL:     g.ImageUrl,
		Genres:       genres,
	}
}
