import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import Underline from '@tiptap/extension-underline'
import Placeholder from '@tiptap/extension-placeholder'
import { useState, useCallback } from 'react'
import './RichEditor.css'

function MenuBar({ editor }) {
  if (!editor) return null

  const addLink = useCallback(() => {
    const url = window.prompt('URL del enlace:')
    if (!url) return
    editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  }, [editor])

  const addImage = useCallback(() => {
    const url = window.prompt('URL de la imagen:')
    if (!url) return
    editor.chain().focus().setImage({ src: url }).run()
  }, [editor])

  return (
    <div className="editor-toolbar">
      <div className="toolbar-group">
        <button type="button" onClick={() => editor.chain().focus().toggleBold().run()}
          className={editor.isActive('bold') ? 'active' : ''} title="Negrita">
          <b>B</b>
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleItalic().run()}
          className={editor.isActive('italic') ? 'active' : ''} title="Cursiva">
          <i>I</i>
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleUnderline().run()}
          className={editor.isActive('underline') ? 'active' : ''} title="Subrayado">
          <u>U</u>
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleStrike().run()}
          className={editor.isActive('strike') ? 'active' : ''} title="Tachado">
          <s>S</s>
        </button>
      </div>

      <div className="toolbar-separator" />

      <div className="toolbar-group">
        <button type="button" onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
          className={editor.isActive('heading', { level: 2 }) ? 'active' : ''} title="Encabezado 2">
          H2
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
          className={editor.isActive('heading', { level: 3 }) ? 'active' : ''} title="Encabezado 3">
          H3
        </button>
        <button type="button" onClick={() => editor.chain().focus().setParagraph().run()}
          className={editor.isActive('paragraph') ? 'active' : ''} title="Párrafo">
          P
        </button>
      </div>

      <div className="toolbar-separator" />

      <div className="toolbar-group">
        <button type="button" onClick={() => editor.chain().focus().toggleBulletList().run()}
          className={editor.isActive('bulletList') ? 'active' : ''} title="Lista">
          &#8226;
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleOrderedList().run()}
          className={editor.isActive('orderedList') ? 'active' : ''} title="Lista numerada">
          1.
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleBlockquote().run()}
          className={editor.isActive('blockquote') ? 'active' : ''} title="Cita">
          &ldquo;
        </button>
        <button type="button" onClick={() => editor.chain().focus().toggleCodeBlock().run()}
          className={editor.isActive('codeBlock') ? 'active' : ''} title="Bloque de código">
          &lt;/&gt;
        </button>
      </div>

      <div className="toolbar-separator" />

      <div className="toolbar-group">
        <button type="button" onClick={addLink}
          className={editor.isActive('link') ? 'active' : ''} title="Enlace">
          🔗
        </button>
        {editor.isActive('link') && (
          <button type="button" onClick={() => editor.chain().focus().unsetLink().run()} title="Quitar enlace">
            ✕
          </button>
        )}
        <button type="button" onClick={addImage} title="Imagen (URL)">
          🖼
        </button>
      </div>

      <div className="toolbar-separator" />

      <div className="toolbar-group">
        <button type="button" onClick={() => editor.chain().focus().setHorizontalRule().run()} title="Línea horizontal">
          ―
        </button>
        <button type="button" onClick={() => editor.chain().focus().undo().run()}
          disabled={!editor.can().undo()} title="Deshacer">
          ↩
        </button>
        <button type="button" onClick={() => editor.chain().focus().redo().run()}
          disabled={!editor.can().redo()} title="Rehacer">
          ↪
        </button>
      </div>
    </div>
  )
}

export default function RichEditor({ content, onChange }) {
  const [showHtml, setShowHtml] = useState(false)
  const [htmlSource, setHtmlSource] = useState(content || '')

  const editor = useEditor({
    extensions: [
      StarterKit,
      Underline,
      Link.configure({
        openOnClick: false,
        HTMLAttributes: { target: '_blank', rel: 'noopener noreferrer' },
      }),
      Image.configure({
        HTMLAttributes: { class: 'article-image' },
      }),
      Placeholder.configure({
        placeholder: 'Escribe el contenido del artículo...',
      }),
    ],
    content: content || '',
    onUpdate: ({ editor }) => {
      const html = editor.getHTML()
      setHtmlSource(html)
      onChange(html)
    },
  })

  const handleHtmlChange = (e) => {
    const val = e.target.value
    setHtmlSource(val)
    if (editor) {
      editor.commands.setContent(val)
    }
    onChange(val)
  }

  const toggleHtmlView = () => {
    if (showHtml) {
      if (editor) editor.commands.setContent(htmlSource)
    } else {
      if (editor) setHtmlSource(editor.getHTML())
    }
    setShowHtml(!showHtml)
  }

  return (
    <div className="rich-editor">
      <div className="editor-header">
        <span className="editor-label">{showHtml ? 'HTML' : 'Editor'}</span>
        <button type="button" className="btn-toggle-html" onClick={toggleHtmlView}>
          {showHtml ? 'Visual' : 'HTML'}
        </button>
      </div>

      {showHtml ? (
        <textarea
          className="html-source"
          value={htmlSource}
          onChange={handleHtmlChange}
          rows={20}
          spellCheck={false}
        />
      ) : (
        <>
          <MenuBar editor={editor} />
          <EditorContent editor={editor} className="editor-content" />
        </>
      )}
    </div>
  )
}
