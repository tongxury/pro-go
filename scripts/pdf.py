# Imports Aspose.Slides for Python via .NET module
import aspose.slides as slides

# Instantiates a Presentation object that represents a presentation file
with slides.Presentation() as presentation:    
    slide = presentation.slides[0]
    slide.shapes.add_auto_shape(slides.ShapeType.LINE, 50, 150, 300, 0)
    presentation.save("NewPresentation.pptx", slides.export.SaveFormat.PPTX)
