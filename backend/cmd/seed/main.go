package main

import (
	"log"

	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found at ../../.env, trying current directory")
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on system environment variables")
		}
	}

	// Initialize Database
	dao.InitDB()

	// Check if events already exist
	var count int64
	dao.DB.Model(&domain.Event{}).Count(&count)
	if count > 0 {
		log.Printf("La base de datos ya contiene %d eventos. Omitiendo la siembra.", count)
		return
	}

	// Define mock events to seed
	mockEvents := []domain.Event{
		{
			Titulo:      "Festival Primavera Sound Buenos Aires",
			Descripcion: "Primavera Sound Buenos Aires reune a artistas consagrados y nuevas promesas en una experiencia pensada para vivir la musica con produccion de primer nivel. El predio cuenta con food trucks, barras tematicas, zonas de descanso y una puesta de luces sincronizada con cada bloque del lineup.",
			Categoria:   "Musica",
			Fecha:       "15 Nov - 20:00 hs",
			HoraInicio:  "20:00",
			HoraFin:     "02:00",
			Ubicacion:   "Parque Sarmiento, CABA",
			Coordenadas: "-34.5678,-58.4890",
			UrlImagen:   "https://images.unsplash.com/photo-1506157786151-b8491531f063?q=80&w=1000",
			Capacidad:   18000,
			Precio:      45000.00,
		},
		{
			Titulo:      "Tech Summit Argentina 2026",
			Descripcion: "Tech Summit esta orientado a estudiantes, profesionales y equipos que quieran ponerse al dia con tendencias de tecnologia, producto e inteligencia artificial. La agenda incluye keynotes, workshops, espacios de networking y una feria de startups con demos en vivo durante toda la jornada.",
			Categoria:   "Charlas",
			Fecha:       "22 Nov - 09:00 hs",
			HoraInicio:  "09:00",
			HoraFin:     "18:30",
			Ubicacion:   "Centro de Convenciones, CABA",
			Coordenadas: "-34.5833,-58.3972",
			UrlImagen:   "https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=1000",
			Capacidad:   3200,
			Precio:      15000.00,
		},
		{
			Titulo:      "Hamlet: Una version contemporanea",
			Descripcion: "Esta puesta de Hamlet combina escenografia minimalista, direccion contemporanea y un elenco con fuerte presencia escenica. La obra propone una lectura actual sobre poder, duelo y decision, con una banda sonora original y recursos multimedia en escena.",
			Categoria:   "Teatro",
			Fecha:       "05 Dic - 21:30 hs",
			HoraInicio:  "21:30",
			HoraFin:     "23:45",
			Ubicacion:   "Teatro San Martin, CABA",
			Coordenadas: "-34.6044,-58.3889",
			UrlImagen:   "https://images.unsplash.com/photo-1507676184212-d03ab07a01bf?q=80&w=1000",
			Capacidad:   800,
			Precio:      8500.00,
		},
		{
			Titulo:      "Final Copa de la Liga",
			Descripcion: "La final de la Copa de la Liga concentra a dos de los equipos mas fuertes del torneo en una noche de estadio repleto y maxima tension competitiva. El operativo de acceso incluye anillos de control, ingresos sectorizados y recomendaciones especiales de llegada para evitar demoras.",
			Categoria:   "Deportes",
			Fecha:       "10 Dic - 17:00 hs",
			HoraInicio:  "17:00",
			HoraFin:     "19:30",
			Ubicacion:   "Estadio Monumental, CABA",
			Coordenadas: "-34.5453,-58.4497",
			UrlImagen:   "https://images.unsplash.com/photo-1508098682722-e99c43a406b2?q=80&w=1000",
			Capacidad:   84000,
			Precio:      25000.00,
		},
	}

	// Insert events
	for _, event := range mockEvents {
		if err := dao.DB.Create(&event).Error; err != nil {
			log.Fatalf("Error al crear evento '%s': %v", event.Titulo, err)
		}
		log.Printf("Evento '%s' creado exitosamente.", event.Titulo)
	}

	log.Println("¡Siembra de eventos completada con éxito!")
}
