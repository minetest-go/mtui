-- contentdb support was removed, keep the entries as manually managed mods
update mod set source_type = 'manual' where source_type = 'cdb';
