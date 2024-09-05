package library

import (
	"context"
	"github.com/carmel/gooxml/document"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"io"
)

// Docx loads text data from an io.Reader.
type Docx struct {
	r io.ReaderAt
	s int64
}

var _ documentloaders.Loader = &Docx{}

func NewDocxLoader(r io.ReaderAt, s int64) *Docx {
	return &Docx{
		r: r,
		s: s,
	}
}

func (o *Docx) Load(_ context.Context) ([]schema.Document, error) {
	doc, err := document.Read(o.r, o.s)
	if err != nil {
		return nil, err
	}

	var docs []schema.Document

	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()

			// add the document to the doc list
			docs = append(docs, schema.Document{
				PageContent: text,
			})
		}
	}

	return docs, nil
}

func (o *Docx) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := o.Load(ctx)
	if err != nil {
		return nil, err
	}

	return textsplitter.SplitDocuments(splitter, docs)
}
