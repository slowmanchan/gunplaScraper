create table gunpla (
  id serial primary key,
  box_art_image_link TEXT,
	title           TEXT,
	subtitle        TEXT,
	classification  TEXT,
	lineup_no        TEXT,
	scale           TEXT,
	franchise       TEXT,
	release_date     TEXT,
	jan_isbn         TEXT,
	run             TEXT,
	price           TEXT,
  includes TEXT,
  features TEXT
)