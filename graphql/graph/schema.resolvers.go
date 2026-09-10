package graph

// Course returns CourseResolver implementation.
func (r *Resolver) Course() CourseResolver { return &courseResolver{r} }

// Professor returns ProfessorResolver implementation.
func (r *Resolver) Professor() ProfessorResolver { return &professorResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Section returns SectionResolver implementation.
func (r *Resolver) Section() SectionResolver { return &sectionResolver{r} }

// SectionWithTime returns SectionWithTimeResolver implementation.
func (r *Resolver) SectionWithTime() SectionWithTimeResolver { return &sectionWithTimeResolver{r} }

type courseResolver struct{ *Resolver }
type professorResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sectionResolver struct{ *Resolver }
type sectionWithTimeResolver struct{ *Resolver }
