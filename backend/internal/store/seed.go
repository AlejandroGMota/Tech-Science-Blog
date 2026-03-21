package store

import (
	"time"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

func SeedArticles(s Store) error {
	articles := []models.Article{
		// Code — Deep Dive
		{
			Slug:        "microservicios-go-produccion",
			Title:       "Lo que aprendí llevando microservicios en Go a producción",
			Excerpt:     "Lecciones reales de migrar un monolito Node.js a microservicios en Go: qué salió bien, qué dolió y qué haría diferente.",
			Category:    "Code",
			ArticleType: "deep_dive",
			Tags:        "Go,Microservicios,Node.js,gRPC,DevOps",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-01"),
			Content:     `<p>Después de 6 meses migrando un monolito Node.js a microservicios en Go, estas son las lecciones que me hubiera gustado saber antes.</p><h2>Por qué Go</h2><p>Concurrencia nativa, binarios estáticos, y un deploy que cabe en un contenedor de 15MB. Para un equipo pequeño, eso es oro.</p><h2>Lo que dolió</h2><p>La falta de generics (en aquel entonces), el manejo de errores verboso, y la tentación de sobre-abstraer. Go te premia por mantenerlo simple.</p><h2>Arquitectura final</h2><p>3 servicios comunicándose por gRPC, un API gateway en Go que expone REST al frontend, y PostgreSQL por servicio. Cada servicio se deploya independientemente.</p>`,
		},
		// Code — Nota Rápida
		{
			Slug:        "error-handling-go-patron",
			Title:       "Un patrón de error handling en Go que me simplificó la vida",
			Excerpt:     "Cómo dejé de repetir if err != nil en cada handler HTTP usando un wrapper simple.",
			Category:    "Code",
			ArticleType: "nota_rapida",
			Tags:        "Go,Error Handling,API,Patrones",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-05"),
			Content:     `<p>Si escribes APIs en Go, probablemente el 40% de tu código es <code>if err != nil</code>. Aquí va un patrón que uso en todos mis proyectos.</p><p>En vez de que cada handler maneje sus propios errores, hago que retornen <code>(any, error)</code> y un middleware se encarga del resto. Menos código, respuestas de error consistentes, y logging centralizado.</p>`,
		},
		// Code — TIL
		{
			Slug:        "til-go-embed",
			Title:       "TIL: go:embed sirve para servir SPAs desde el binario",
			Excerpt:     "Hoy descubrí que puedo embeber el build de React directo en el binario de Go.",
			Category:    "Code",
			ArticleType: "til",
			Tags:        "Go,React,Deploy,Embed",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-10"),
			Content:     `<p>Hoy descubrí <code>//go:embed</code> y cambió cómo deployeo este blog. En vez de copiar archivos estáticos al servidor, el binario de Go ya los trae dentro. Un solo archivo, zero config de Nginx para el frontend.</p>`,
		},
		// Business — Deep Dive
		{
			Slug:        "importar-china-mexico-lecciones",
			Title:       "Lo que nadie te dice de importar desde China a México",
			Excerpt:     "Aduanas, agentes, tiempos muertos y el SAT: todo lo que aprendí montando la cadena de suministro de Celinki.",
			Category:    "Business",
			ArticleType: "deep_dive",
			Tags:        "Importación,China,México,Celinki,Aduanas,SAT",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-02-15"),
			Content:     `<p>Cuando empecé Celinki, pensé que importar desde China era cuestión de encontrar un proveedor en Alibaba y pagar. Estaba equivocado por un factor de 10x en complejidad.</p><h2>El agente aduanal</h2><p>Es tu mejor amigo o tu peor pesadilla. Un buen agente te ahorra semanas y miles de pesos. Un malo te puede dejar mercancía atorada en aduana por meses.</p><h2>Tiempos reales</h2><p>Del pedido a tu bodega: 45-90 días si todo sale bien. Agrega 2-3 semanas si el SAT decide revisar tu contenedor.</p><h2>Lo que haría diferente</h2><p>Empezar con pedidos pequeños por courier (DHL/FedEx) antes de comprometerse con un contenedor completo. La señal de demanda vale más que el ahorro por volumen.</p>`,
		},
		// Business — Nota Rápida
		{
			Slug:        "sat-facturacion-tips",
			Title:       "3 cosas del SAT que todo emprendedor debería saber desde el día 1",
			Excerpt:     "Errores de facturación que me costaron dinero y cómo evitarlos.",
			Category:    "Business",
			ArticleType: "nota_rapida",
			Tags:        "SAT,Facturación,Emprendimiento,Fiscal",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-02-20"),
			Content:     `<p>1. Tu régimen fiscal importa más de lo que crees — cambiarlo después es un dolor. 2. Factura TODO desde el día 1, aunque sea incómodo. 3. Guarda tus XMLs, no solo los PDFs. El SAT pide XMLs y si los pierdes, el contador no puede hacer magia.</p>`,
		},
		// Ideas — Deep Dive
		{
			Slug:        "semiconductores-geopolitica-2025",
			Title:       "Por qué los semiconductores son el nuevo petróleo geopolítico",
			Excerpt:     "TSMC, ASML y la guerra silenciosa entre EE.UU. y China por el control de los chips más avanzados del mundo.",
			Category:    "Ideas",
			ArticleType: "deep_dive",
			Tags:        "Geopolítica,Semiconductores,TSMC,China,Nearshoring",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-01-20"),
			Content:     `<p>Si quieres entender la geopolítica actual, sigue los chips. No los de póker — los semiconductores.</p><h2>El cuello de botella</h2><p>TSMC en Taiwán fabrica el 90% de los chips más avanzados del mundo. ASML en Holanda hace las únicas máquinas capaces de producirlos. Dos empresas, dos países, y el mundo entero depende de ellos.</p><h2>La jugada de EE.UU.</h2><p>Las restricciones de exportación a China no son solo comerciales — son una estrategia para mantener una ventaja tecnológica de al menos una generación. China lo sabe y está invirtiendo $150B en independencia de chips.</p><h2>¿Y México?</h2><p>Nearshoring de semiconductores es la oportunidad más grande que ha tenido México en décadas. Pero necesita infraestructura, talento y estabilidad regulatoria para capturarla.</p>`,
		},
		// Ideas — Nota Rápida
		{
			Slug:        "opinion-ai-reemplazar-programadores",
			Title:       "No, la IA no va a reemplazar programadores (pero sí va a cambiar qué significa programar)",
			Excerpt:     "Mi opinión sobre el debate más repetido de 2025.",
			Category:    "Ideas",
			ArticleType: "nota_rapida",
			Tags:        "IA,Programación,Opinión,Futuro del trabajo",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-12"),
			Content:     `<p>Cada vez que alguien dice "la IA va a reemplazar a los programadores", lo que realmente está pasando es que la IA está eliminando el código boilerplate. Y eso es bueno — nadie debería pasar 3 horas configurando un CRUD.</p><p>Lo que no reemplaza: entender el problema del usuario, diseñar arquitectura que escale, y tomar decisiones de trade-offs. Eso sigue siendo humano. Por ahora.</p>`,
		},
		// Stack de vida — Deep Dive
		{
			Slug:        "sistema-productividad-dev-emprendedor",
			Title:       "Mi sistema de productividad como dev y emprendedor",
			Excerpt:     "Cómo organizo mi tiempo entre código, Celinki y aprendizaje sin quemarme.",
			Category:    "Stack de vida",
			ArticleType: "deep_dive",
			Tags:        "Productividad,Notion,Time Blocking,Emprendimiento",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-02-01"),
			Content:     `<p>Hacer malabares entre un trabajo técnico, un negocio de importación y proyectos personales requiere un sistema. No disciplina heroica — un sistema.</p><h2>Time blocking</h2><p>Mañanas para código (deep work). Tardes para Celinki (operaciones, proveedores). Noches para aprender algo nuevo o escribir. Los fines de semana son sagrados.</p><h2>Herramientas</h2><p>Notion para proyectos largos, Todoist para el día a día, y un cuaderno físico para pensar. La herramienta importa menos que la consistencia.</p><h2>La regla del 80%</h2><p>Si esperas a tener todo perfecto, nunca publicas, nunca lanzas, nunca avanzas. Hecho al 80% y publicado vale más que perfecto en drafts.</p>`,
		},
		// Stack de vida — TIL
		{
			Slug:        "til-terminal-multiplexer",
			Title:       "TIL: tmux + ssh = nunca más perder una sesión remota",
			Excerpt:     "Descubrí tmux y ahora mis sesiones SSH sobreviven a desconexiones.",
			Category:    "Stack de vida",
			ArticleType: "til",
			Tags:        "tmux,SSH,Terminal,Linux,Herramientas",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-08"),
			Content:     `<p>Llevaba meses perdiendo trabajo cuando se caía la conexión SSH a mi VM de Oracle Cloud. Hoy descubrí tmux: <code>tmux new -s work</code>, y si se cae la conexión, <code>tmux attach -t work</code> y todo sigue ahí. Game changer.</p>`,
		},
		// Aprendiendo en público — Deep Dive
		{
			Slug:        "aprendiendo-oracle-cloud-free-tier",
			Title:       "Monté un blog en Oracle Cloud Free Tier y esto es lo que aprendí",
			Excerpt:     "De no saber nada de Oracle a tener un servidor ARM con Nginx, Go y una base de datos autónoma — gratis.",
			Category:    "Aprendiendo en público",
			ArticleType: "deep_dive",
			Tags:        "Oracle Cloud,ARM,Nginx,CI/CD,Infraestructura",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-01-10"),
			Content:     `<p>Elegí Oracle Cloud por una razón: el free tier incluye una VM ARM con 4 OCPUs y 24GB de RAM. Eso es más que suficiente para un blog y varios side projects.</p><h2>Lo bueno</h2><p>El hardware es excelente para el precio (gratis). La base de datos autónoma te ahorra administración. El networking es configurable.</p><h2>Lo malo</h2><p>La documentación es un laberinto. La consola web es lenta. Y configurar el firewall requiere cambios tanto en la consola OCI como en iptables de la VM.</p><h2>Lo que aprendí</h2><p>Systemd para mantener servicios corriendo, Nginx como reverse proxy, Let's Encrypt para SSL, y GitHub Actions para CI/CD. Todo desde cero.</p>`,
		},
		// Aprendiendo en público — TIL
		{
			Slug:        "til-css-container-queries",
			Title:       "TIL: CSS Container Queries cambian todo para componentes responsive",
			Excerpt:     "Media queries miran el viewport. Container queries miran el contenedor. Eso lo cambia todo.",
			Category:    "Aprendiendo en público",
			ArticleType: "til",
			Tags:        "CSS,Frontend,Responsive,Container Queries",
			Author:      "Alejandro G. Mota",
			PublishedAt: parseDate("2025-03-15"),
			Content:     `<p>Hoy aprendí que CSS tiene container queries (<code>@container</code>) y resuelven un problema que llevo años parchando con media queries. En vez de preguntar "¿qué tan ancho es el viewport?", preguntan "¿qué tan ancho es mi contenedor?". Perfecto para componentes reutilizables que viven en layouts diferentes.</p>`,
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
