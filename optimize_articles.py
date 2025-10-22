#!/usr/bin/env python3
"""
Script to optimize all blog articles with comprehensive SEO improvements
"""

import json
import re
from pathlib import Path

# Article metadata from entradas.json
articles_data = {
    "avances-más-recientes-en-web-3.html": {
        "title": "Avances Más Recientes en Web 3",
        "date": "2024-12-04",
        "category": "Blockchain",
        "description": "La Web 3, también conocida como la web descentralizada, está transformando la manera en que interactuamos con la tecnología. Descubre los avances más recientes en blockchain, dApps y contratos inteligentes.",
        "keywords": "web 3, blockchain, descentralización, contratos inteligentes, dApps, ethereum 2.0, DeFi, NFT",
        "section": "Blockchain",
        "tags": ["Web 3", "Blockchain", "Descentralización", "Tecnología"]
    },
    "como-chat-gpt-cambiará-tu-vida.html": {
        "title": "Cómo Chat GPT Cambiará Tu Vida",
        "date": "2024-12-05",
        "category": "IA",
        "description": "La inteligencia artificial ha llegado para revolucionar la forma en que interactuamos con el mundo. Descubre cómo ChatGPT está transformando la productividad, educación y comunicación.",
        "keywords": "ChatGPT, inteligencia artificial, IA conversacional, transformación digital, productividad, educación",
        "section": "Inteligencia Artificial",
        "tags": ["ChatGPT", "IA", "Productividad", "Educación"]
    },
    "cómo-la-inteligencia-artificial-está-transformando-tu-día-a-día.html": {
        "title": "Cómo la Inteligencia Artificial Está Transformando Tu Día a Día",
        "date": "2024-12-09",
        "category": "IA",
        "description": "La inteligencia artificial está cambiando la manera en que vivimos, trabajamos y nos relacionamos. Explora las aplicaciones de IA en la vida cotidiana, desde asistentes virtuales hasta recomendaciones personalizadas.",
        "keywords": "inteligencia artificial, IA en la vida cotidiana, asistentes virtuales, machine learning, automatización, tecnología AI",
        "section": "Inteligencia Artificial",
        "tags": ["IA", "Vida Cotidiana", "Automatización", "Tecnología"]
    },
    "escombro-en-la-agronomía-y-tecnología.html": {
        "title": "Escombro en la Agronomía y Tecnología",
        "date": "2024-12-04",
        "category": "IA",
        "description": "El aprovechamiento de residuos de construcción y demolición ha surgido como una alternativa sostenible en agronomía. Descubre cómo la tecnología está revolucionando el uso de escombros en la agricultura.",
        "keywords": "agronomía sostenible, reciclaje de escombros, agricultura tecnológica, sostenibilidad, residuos de construcción",
        "section": "Agronomía y Tecnología",
        "tags": ["Agronomía", "Sostenibilidad", "Reciclaje", "Tecnología"]
    },
    "la-mejor-forma-de-dormir.html": {
        "title": "La Mejor Forma de Dormir",
        "date": "2024-12-05",
        "category": "IA",
        "description": "Dormir es una de las actividades más importantes para el bienestar general. Descubre las mejores prácticas y posiciones para un sueño reparador basadas en la ciencia.",
        "keywords": "calidad del sueño, cómo dormir mejor, higiene del sueño, salud, bienestar, ciencia del sueño",
        "section": "Salud y Bienestar",
        "tags": ["Salud", "Bienestar", "Sueño", "Ciencia"]
    },
    "por-qué-no-debes-usar-generadores-de-imágenes-como-estrategia-de-contenido.html": {
        "title": "Por Qué No Debes Usar Generadores de Imágenes Como Estrategia de Contenido",
        "date": "2024-12-06",
        "category": "IA",
        "description": "En la era de la inteligencia artificial, los generadores de imágenes han ganado popularidad. Sin embargo, descubre por qué no son la mejor estrategia para tu contenido digital.",
        "keywords": "generadores de imágenes IA, estrategia de contenido, marketing digital, IA generativa, contenido visual, autenticidad",
        "section": "Marketing Digital",
        "tags": ["IA", "Marketing", "Contenido Digital", "Estrategia"]
    },
    "sensores-agrónomos-con-ia-y-tecnología.html": {
        "title": "Sensores Agrónomos con IA y Tecnología",
        "date": "2024-12-03",
        "category": "IA",
        "description": "La inteligencia artificial está marcando un cambio significativo en la agricultura. Descubre cómo los sensores agrónomos con IA están revolucionando la gestión de cultivos.",
        "keywords": "agricultura inteligente, sensores agrónomos, IA en agricultura, tecnología agrícola, agricultura de precisión, smart farming",
        "section": "Agronomía y Tecnología",
        "tags": ["Agronomía", "IA", "Sensores", "Tecnología"]
    },
    "sistemas-de-riego-entre-países-desarrollados-y-menos-desarrollados.html": {
        "title": "Sistemas de Riego Entre Países Desarrollados y Menos Desarrollados",
        "date": "2024-12-03",
        "category": "Agronomía",
        "description": "El uso de drones para optimizar los sistemas de riego está revolucionando la agricultura mundial. Compara las tecnologías de riego entre países desarrollados y en desarrollo.",
        "keywords": "sistemas de riego, drones agrícolas, agricultura tecnológica, países desarrollados, agricultura sostenible, innovación agrícola",
        "section": "Agronomía y Tecnología",
        "tags": ["Agronomía", "Drones", "Tecnología", "Sostenibilidad"]
    },
    "usos-del-data-science-en-la-confección-de-ropa.html": {
        "title": "Usos del Data Science en la Confección de Ropa",
        "date": "2024-12-05",
        "category": "IA",
        "description": "La industria de la confección está experimentando una transformación radical gracias al Data Science. Explora cómo el análisis de datos está revolucionando el diseño y producción de ropa.",
        "keywords": "data science, industria textil, confección de ropa, análisis de datos, moda tecnológica, industria 4.0",
        "section": "Data Science",
        "tags": ["Data Science", "Moda", "Tecnología", "Industria"]
    }
}

def create_seo_head(article_data, filename):
    """Generate comprehensive SEO head section"""
    url_filename = filename.replace(' ', '-').lower()

    head_template = f'''<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <!-- Title -->
    <title>{article_data['title']} | Tech & Science Blog</title>

    <!-- Meta Description -->
    <meta name="description" content="{article_data['description']}">

    <!-- Meta Keywords -->
    <meta name="keywords" content="{article_data['keywords']}">

    <!-- Meta Tags -->
    <meta name="author" content="Tech & Science Blog">
    <meta name="robots" content="index, follow, max-image-preview:large, max-snippet:-1">
    <meta name="date" content="{article_data['date']}">

    <!-- Canonical URL -->
    <link rel="canonical" href="https://blog.alejandrogmota.com/entradas/{url_filename}">

    <!-- Open Graph Meta Tags -->
    <meta property="og:locale" content="es_ES">
    <meta property="og:type" content="article">
    <meta property="og:title" content="{article_data['title']}">
    <meta property="og:description" content="{article_data['description']}">
    <meta property="og:image" content="https://blog.alejandrogmota.com/img/logo.webp">
    <meta property="og:image:width" content="1200">
    <meta property="og:image:height" content="630">
    <meta property="og:url" content="https://blog.alejandrogmota.com/entradas/{url_filename}">
    <meta property="og:site_name" content="Tech & Science Blog">
    <meta property="article:published_time" content="{article_data['date']}T00:00:00+00:00">
    <meta property="article:author" content="Tech & Science Blog">
    <meta property="article:section" content="{article_data['section']}">'''

    for tag in article_data['tags']:
        head_template += f'\n    <meta property="article:tag" content="{tag}">'

    head_template += f'''

    <!-- Twitter Meta Tags -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="{article_data['title']}">
    <meta name="twitter:description" content="{article_data['description'][:150]}">
    <meta name="twitter:image" content="https://blog.alejandrogmota.com/img/logo.webp">
    <meta name="twitter:site" content="@TechScienceBlog">
    <meta name="twitter:creator" content="@TechScienceBlog">

    <!-- Google Adsense -->
    <meta name="google-adsense-account" content="ca-pub-2557598346205228">

    <!-- Preconnect for Performance -->
    <link rel="preconnect" href="https://blog.alejandrogmota.com">
    <link rel="dns-prefetch" href="https://via.placeholder.com">

    <!-- Favicon -->
    <link rel="icon" href="https://blog.alejandrogmota.com/favicon.PNG" type="image/x-icon">
    <link rel="apple-touch-icon" href="https://blog.alejandrogmota.com/img/logo.webp">

    <!-- Stylesheet -->
    <link rel="stylesheet" href="../css/styles.css">

    <!-- Web Manifest -->
    <link rel="manifest" href="/manifest.json">

    <!-- Theme Color -->
    <meta name="theme-color" content="#0044cc">

    <!-- Structured Data - Article -->
    <script type="application/ld+json">
    {{
      "@context": "https://schema.org",
      "@type": "BlogPosting",
      "headline": "{article_data['title']}",
      "description": "{article_data['description']}",
      "image": "https://blog.alejandrogmota.com/img/logo.webp",
      "author": {{
        "@type": "Organization",
        "name": "Tech & Science Blog"
      }},
      "publisher": {{
        "@type": "Organization",
        "name": "Tech & Science Blog",
        "logo": {{
          "@type": "ImageObject",
          "url": "https://blog.alejandrogmota.com/img/logo.webp"
        }}
      }},
      "datePublished": "{article_data['date']}",
      "dateModified": "{article_data['date']}",
      "mainEntityOfPage": {{
        "@type": "WebPage",
        "@id": "https://blog.alejandrogmota.com/entradas/{url_filename}"
      }},
      "articleSection": "{article_data['section']}",
      "keywords": {json.dumps(article_data['keywords'].split(', '))},
      "inLanguage": "es-ES",
      "isPartOf": {{
        "@type": "Blog",
        "name": "Tech & Science Blog",
        "url": "https://blog.alejandrogmota.com"
      }}
    }}
    </script>

    <!-- Breadcrumb Structured Data -->
    <script type="application/ld+json">
    {{
      "@context": "https://schema.org",
      "@type": "BreadcrumbList",
      "itemListElement": [{{
        "@type": "ListItem",
        "position": 1,
        "name": "Inicio",
        "item": "https://blog.alejandrogmota.com/"
      }},{{
        "@type": "ListItem",
        "position": 2,
        "name": "Entradas",
        "item": "https://blog.alejandrogmota.com/html/entradas.html"
      }},{{
        "@type": "ListItem",
        "position": 3,
        "name": "{article_data['title']}",
        "item": "https://blog.alejandrogmota.com/entradas/{url_filename}"
      }}]
    }}
    </script>
</head>'''

    return head_template

def create_breadcrumb(title):
    """Generate breadcrumb navigation"""
    return f'''
        <!-- Breadcrumbs Navigation -->
        <nav aria-label="breadcrumb" style="padding: 20px 10%; background-color: #f9f9f9;">
            <ol style="list-style: none; display: flex; gap: 10px; padding: 0; margin: 0; font-size: 0.9rem;">
                <li><a href="https://blog.alejandrogmota.com/" style="color: #0044cc; text-decoration: none;">Inicio</a></li>
                <li style="color: #666;">/</li>
                <li><a href="https://blog.alejandrogmota.com/html/entradas.html" style="color: #0044cc; text-decoration: none;">Entradas</a></li>
                <li style="color: #666;">/</li>
                <li style="color: #333; font-weight: bold;">{title}</li>
            </ol>
        </nav>
'''

def optimize_article(filepath, article_data):
    """Optimize a single article with SEO improvements"""
    print(f"Optimizing {filepath.name}...")

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Extract body content (everything after </head>)
    body_match = re.search(r'</head>\s*<body>(.*)</body>', content, re.DOTALL)
    if not body_match:
        print(f"  ⚠ Could not extract body from {filepath.name}")
        return False

    body_content = body_match.group(1)

    # Generate new head
    new_head = create_seo_head(article_data, filepath.name)

    # Add breadcrumbs after header
    breadcrumb = create_breadcrumb(article_data['title'])

    # Replace <h2> with <h1> for main title and add semantic time tag
    body_content = re.sub(
        r'<h2>([^<]+)</h2>\s*<p><strong>Fecha de Publicación:</strong> (\d{4}-\d{2}-\d{2})</p>',
        r'<h1>\1</h1>\n            <p><strong>Fecha de Publicación:</strong> <time datetime="\2">\2</time></p>',
        body_content
    )

    # Add breadcrumbs after header
    body_content = re.sub(
        r'(</header>)',
        r'\1' + breadcrumb,
        body_content,
        count=1
    )

    # Optimize images (add lazy loading, dimensions, and better alt text)
    body_content = re.sub(
        r'<img src="([^"]+)" alt="([^"]+)">',
        r'<img src="\1" alt="\2" loading="lazy" width="800" height="400" style="max-width: 100%; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">',
        body_content
    )

    # Reconstruct full HTML
    new_content = new_head + '\n    <body>' + body_content + '</body>\n    </html>\n    '

    # Write optimized content
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(new_content)

    print(f"  ✓ {filepath.name} optimized successfully")
    return True

def main():
    """Main function to optimize all articles"""
    entradas_dir = Path('/home/user/Tech-Science-Blog/entradas')

    print("Starting SEO optimization for all articles...")
    print("=" * 60)

    optimized = 0
    skipped = 0

    for filename, article_data in articles_data.items():
        filepath = entradas_dir / filename

        # Skip agentes-de-ia.html as it's already optimized
        if filename == 'agentes-de-ia.html':
            print(f"Skipping {filename} (already optimized)")
            skipped += 1
            continue

        if filepath.exists():
            if optimize_article(filepath, article_data):
                optimized += 1
        else:
            print(f"  ⚠ File not found: {filename}")
            skipped += 1

    print("=" * 60)
    print(f"Optimization complete!")
    print(f"  ✓ Optimized: {optimized} articles")
    print(f"  ⚠ Skipped: {skipped} articles")

if __name__ == "__main__":
    main()
