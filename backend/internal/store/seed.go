package store

import (
	"time"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

func SeedArticles(s Store) error {
	articles := []models.Article{
		{
			Slug:        "agentes-de-ia",
			Title:       "Agentes De IA",
			Excerpt:     "Los agentes de inteligencia artificial (IA) están revolucionando la manera en que interactuamos con la tecnología.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-03"),
			Content:     `<p>Los agentes de inteligencia artificial (IA) están revolucionando la manera en que interactuamos con la tecnología. Estos sistemas autónomos pueden tomar decisiones, aprender de su entorno y ejecutar tareas complejas sin intervención humana directa.</p><h2>¿Qué son los agentes de IA?</h2><p>Un agente de IA es un sistema que percibe su entorno a través de sensores y actúa sobre él mediante actuadores. Pueden ser tan simples como un termostato inteligente o tan complejos como un vehículo autónomo.</p><h2>Tipos de agentes</h2><p>Existen varios tipos: agentes reactivos simples, agentes basados en modelos, agentes basados en objetivos y agentes basados en utilidad. Cada uno tiene diferentes niveles de complejidad y capacidad de decisión.</p>`,
		},
		{
			Slug:        "avances-mas-recientes-en-web-3",
			Title:       "Avances Más Recientes En Web 3",
			Excerpt:     "La Web 3, también conocida como la web descentralizada, está transformando la manera en que interactuamos con la tecnología.",
			Category:    "Blockchain",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-04"),
			Content:     `<p>La Web 3, también conocida como la web descentralizada, está transformando la manera en que interactuamos con la tecnología. Basada en blockchain, ofrece mayor transparencia, seguridad y control para los usuarios.</p><h2>Principales avances</h2><p>Entre los avances más recientes se encuentran los contratos inteligentes más eficientes, las finanzas descentralizadas (DeFi) y los tokens no fungibles (NFTs) con aplicaciones prácticas más allá del arte digital.</p>`,
		},
		{
			Slug:        "como-chat-gpt-cambiara-tu-vida",
			Title:       "Cómo Chat GPT Cambiará Tu Vida",
			Excerpt:     "La inteligencia artificial (IA) ha llegado para revolucionar la forma en que interactuamos con el mundo.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-05"),
			Content:     `<p>La inteligencia artificial (IA) ha llegado para revolucionar la forma en que interactuamos con el mundo. ChatGPT, desarrollado por OpenAI, es uno de los modelos de lenguaje más avanzados y accesibles.</p><h2>Impacto en el trabajo</h2><p>Desde la redacción de correos hasta la programación, ChatGPT está transformando la productividad en múltiples industrias.</p><h2>Educación</h2><p>Los estudiantes pueden usar esta herramienta como tutor personalizado, obteniendo explicaciones adaptadas a su nivel de conocimiento.</p>`,
		},
		{
			Slug:        "como-la-inteligencia-artificial-esta-transformando-tu-dia-a-dia",
			Title:       "Cómo La Inteligencia Artificial Está Transformando Tu Día A Día",
			Excerpt:     "La inteligencia artificial (IA) está cambiando la manera en que vivimos, trabajamos y nos relacionamos.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-09"),
			Content:     `<p>La inteligencia artificial (IA) está cambiando la manera en que vivimos, trabajamos y nos relacionamos. Desde asistentes virtuales hasta recomendaciones personalizadas, la IA se ha integrado en nuestra rutina diaria.</p><h2>En el hogar</h2><p>Dispositivos inteligentes, termostatos que aprenden tus preferencias y sistemas de seguridad basados en reconocimiento facial son solo algunos ejemplos.</p><h2>En la salud</h2><p>La IA está revolucionando el diagnóstico médico, permitiendo detectar enfermedades de forma temprana con mayor precisión.</p>`,
		},
		{
			Slug:        "escombro-en-la-agronomia-y-tecnologia",
			Title:       "Escombro En La Agronomía Y Tecnología",
			Excerpt:     "El aprovechamiento de residuos de construcción y demolición ha surgido como una alternativa sostenible.",
			Category:    "Agronomía",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-04"),
			Content:     `<p>El aprovechamiento de residuos de construcción y demolición, o escombros, ha surgido como una alternativa sostenible en un mundo que busca soluciones para reducir el impacto ambiental.</p><h2>Aplicaciones en agricultura</h2><p>Los materiales reciclados pueden utilizarse para mejorar la estructura del suelo, crear sistemas de drenaje y como base para invernaderos sostenibles.</p>`,
		},
		{
			Slug:        "la-mejor-forma-de-dormir",
			Title:       "La Mejor Forma De Dormir",
			Excerpt:     "Dormir es una de las actividades más importantes para el bienestar general de las personas.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-05"),
			Content:     `<p>Dormir es una de las actividades más importantes para el bienestar general de las personas. La tecnología moderna nos ayuda a entender mejor nuestros patrones de sueño.</p><h2>Tecnología del sueño</h2><p>Dispositivos wearables y aplicaciones de IA pueden monitorear las fases del sueño y ofrecer recomendaciones personalizadas para mejorar la calidad del descanso.</p>`,
		},
		{
			Slug:        "por-que-no-debes-usar-generadores-de-imagenes",
			Title:       "Por Qué No Debes Usar Generadores De Imágenes Como Estrategia De Contenido",
			Excerpt:     "Los generadores de imágenes han ganado popularidad como herramienta rápida y económica para la creación de contenido visual.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-06"),
			Content:     `<p>En la era de la inteligencia artificial, los generadores de imágenes han ganado popularidad como una herramienta rápida y económica para la creación de contenido visual. Sin embargo, depender exclusivamente de ellos puede ser contraproducente.</p><h2>Problemas de autenticidad</h2><p>Las imágenes generadas por IA carecen de la autenticidad que los usuarios buscan en las marcas y el contenido digital.</p><h2>Implicaciones legales</h2><p>El uso de imágenes generadas por IA plantea cuestiones sobre derechos de autor y propiedad intelectual que aún no están completamente resueltas.</p>`,
		},
		{
			Slug:        "sensores-agronomos-con-ia-y-tecnologia",
			Title:       "Sensores Agrónomos Con IA Y Tecnología",
			Excerpt:     "La inteligencia artificial está marcando un cambio significativo en la gestión de cultivos agrícolas.",
			Category:    "IA",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-03"),
			Content:     `<p>La inteligencia artificial (IA) está marcando un cambio significativo en la manera de gestionar los cultivos agrícolas. Los sensores inteligentes permiten monitorear en tiempo real las condiciones del suelo, el clima y la salud de las plantas.</p><h2>Tipos de sensores</h2><p>Sensores de humedad, temperatura, pH y nutrientes del suelo se combinan con algoritmos de IA para optimizar el uso de recursos.</p>`,
		},
		{
			Slug:        "sistemas-de-riego-entre-paises",
			Title:       "Sistemas De Riego Entre Países Desarrollados Y Menos Desarrollados",
			Excerpt:     "El uso de drones para optimizar los sistemas de riego está revolucionando la agricultura mundial.",
			Category:    "Agronomía",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-03"),
			Content:     `<p>El uso de drones para optimizar los sistemas de riego está revolucionando la agricultura mundial. Sin embargo, existe una brecha significativa entre los países desarrollados y los menos desarrollados en la adopción de estas tecnologías.</p><h2>Tecnología de riego avanzada</h2><p>Los sistemas de riego por goteo controlados por IA pueden reducir el consumo de agua hasta en un 60% comparado con métodos tradicionales.</p>`,
		},
		{
			Slug:        "usos-del-data-science-en-la-confeccion-de-ropa",
			Title:       "Usos Del Data Science En La Confección De Ropa",
			Excerpt:     "La industria de la confección está experimentando una transformación radical gracias al Data Science.",
			Category:    "Data Science",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2024-12-05"),
			Content:     `<p>La industria de la confección está experimentando una transformación radical gracias al impacto del Data Science. Desde la predicción de tendencias hasta la optimización de la cadena de suministro.</p><h2>Predicción de tendencias</h2><p>Algoritmos de machine learning analizan datos de redes sociales, pasarelas y ventas para predecir qué estilos serán populares en las próximas temporadas.</p><h2>Personalización</h2><p>El data science permite crear experiencias de compra personalizadas, recomendando tallas, estilos y colores basados en el historial del usuario.</p>`,
		},
	}

	for i := range articles {
		if err := s.CreateArticle(&articles[i]); err != nil {
			// Skip if already exists
			continue
		}
	}
	return nil
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
