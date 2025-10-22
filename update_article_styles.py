#!/usr/bin/env python3
"""
Script to update article styles to use CSS classes instead of inline styles
"""

import re
from pathlib import Path

def update_breadcrumb_styles(content):
    """Remove inline styles from breadcrumbs"""
    # Pattern to match the old breadcrumb style
    old_breadcrumb = r'<nav aria-label="breadcrumb" style="[^"]*">\s*<ol style="[^"]*">'
    new_breadcrumb = '<nav aria-label="breadcrumb">\n            <ol>'
    content = re.sub(old_breadcrumb, new_breadcrumb, content)

    # Remove inline styles from breadcrumb links
    content = re.sub(
        r'<li><a href="([^"]+)" style="[^"]*">([^<]+)</a></li>',
        r'<li><a href="\1">\2</a></li>',
        content
    )

    return content

def update_h3_styles(content):
    """Remove inline styles from h3 tags"""
    content = re.sub(
        r'<h3 style="[^"]*">',
        '<h3>',
        content
    )
    return content

def update_image_styles(content):
    """Remove inline styles from images and their containers"""
    # Update image divs
    content = re.sub(
        r'<div class="about-image" style="[^"]*">',
        '<div class="about-image">',
        content
    )

    # Update images - keep lazy loading and dimensions, remove inline styles
    content = re.sub(
        r'<img src="([^"]+)" alt="([^"]+)"([^>]*) style="[^"]*">',
        r'<img src="\1" alt="\2"\3>',
        content
    )

    return content

def update_article(filepath):
    """Update a single article file"""
    print(f"Updating {filepath.name}...")

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        # Apply updates
        content = update_breadcrumb_styles(content)
        content = update_h3_styles(content)
        content = update_image_styles(content)

        # Write back
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)

        print(f"  ✓ {filepath.name} updated successfully")
        return True

    except Exception as e:
        print(f"  ✗ Error updating {filepath.name}: {e}")
        return False

def main():
    """Main function"""
    entradas_dir = Path('/home/user/Tech-Science-Blog/entradas')

    print("Updating article styles to use CSS classes...")
    print("=" * 60)

    # Get all HTML files
    html_files = sorted(entradas_dir.glob('*.html'))

    updated = 0
    failed = 0

    for filepath in html_files:
        if update_article(filepath):
            updated += 1
        else:
            failed += 1

    print("=" * 60)
    print(f"Update complete!")
    print(f"  ✓ Updated: {updated} articles")
    if failed > 0:
        print(f"  ✗ Failed: {failed} articles")

if __name__ == "__main__":
    main()
