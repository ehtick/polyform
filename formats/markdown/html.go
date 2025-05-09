package markdown

import "github.com/EliCDavis/polyform/formats/txt"

type htmlWriter struct {
	writer *txt.Writer
}

func (w *htmlWriter) openTag(element []byte) {
	w.writer.Append([]byte("<"))
	w.writer.Append(element)
	w.writer.Append([]byte(">"))
}

func (w *htmlWriter) closeTag(element []byte) {
	w.writer.Append([]byte("</"))
	w.writer.Append(element)
	w.writer.Append([]byte(">\n"))
}

func (w *htmlWriter) element(element []byte, text string) (int, error) {
	w.writer.StartEntry()
	w.openTag(element)
	w.writer.String(text)
	w.closeTag(element)
	return w.writer.FinishEntry()
}

func (w *htmlWriter) elementWithId(element []byte, text string, id string) (int, error) {
	w.writer.StartEntry()
	w.writer.Append([]byte("<"))
	w.writer.Append(element)
	w.writer.Append([]byte(" id=\""))
	w.writer.String(id)
	w.writer.Append([]byte("\">"))
	w.writer.String(text)
	w.closeTag(element)
	return w.writer.FinishEntry()
}

func (w *htmlWriter) Header1(text string) (int, error) {
	return w.element([]byte("h1"), text)
}

func (w *htmlWriter) Header2(text string) (int, error) {
	return w.element([]byte("h2"), text)
}

func (w *htmlWriter) Header2WithId(text string, id string) (int, error) {
	return w.elementWithId([]byte("h2"), text, id)
}

func (w *htmlWriter) Header3(text string) (int, error) {
	return w.element([]byte("h3"), text)
}

func (w *htmlWriter) Header3WithId(text string, id string) (int, error) {
	return w.elementWithId([]byte("h3"), text, id)
}

func (w *htmlWriter) Link(text string, link string) (int, error) {
	w.writer.StartEntry()
	w.writer.Append([]byte("<a href=\"#"))
	w.writer.String(link)
	w.writer.Append([]byte("\">"))
	w.writer.String(text)
	w.closeTag([]byte("a"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) Paragraph(text string) (int, error) {
	return w.element([]byte("p"), text)
}

func (w *htmlWriter) NewLine() (int, error) {
	w.writer.StartEntry()
	w.openTag([]byte("br"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) StartBulletList() (int, error) {
	w.writer.StartEntry()
	w.openTag([]byte("ul"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) EndBulletList() (int, error) {
	w.writer.StartEntry()
	w.closeTag([]byte("ul"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) StartBullet() (int, error) {
	w.writer.StartEntry()
	w.openTag([]byte("li"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) EndBullet() (int, error) {
	w.writer.StartEntry()
	w.closeTag([]byte("li"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) StartBold() (int, error) {
	w.writer.StartEntry()
	w.openTag([]byte("b"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) EndBold() (int, error) {
	w.writer.StartEntry()
	w.closeTag([]byte("b"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) StartItalics() (int, error) {
	w.writer.StartEntry()
	w.openTag([]byte("i"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) EndItalics() (int, error) {
	w.writer.StartEntry()
	w.closeTag([]byte("i"))
	return w.writer.FinishEntry()
}

func (w *htmlWriter) Text(text string) (int, error) {
	w.writer.StartEntry()
	w.writer.String(text)
	return w.writer.FinishEntry()
}

func (w *htmlWriter) Error() error {
	return w.writer.Error()
}
