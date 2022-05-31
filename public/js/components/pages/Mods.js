export default {
	template: /*html*/`
		<div>
			<h3>Mod selection</h3>
			<table class="table table-condensed">
				<thead>
					<tr>
						<th>Type</th>
						<th>Name</th>
						<th>Source-Type</th>
						<th>Source</th>
						<th>Version</th>
						<th>Auto-update</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>
							<span class="badge badge-secondary">Mod</span>
						</td>
						<td>moreblocks</td>
						<td>
							<span class="badge badge-success">
								Git
							</span>
						</td>
						<td>
							<a href="https://github.com/minetest-mods/moreblocks.git">
								https://github.com/minetest-mods/moreblocks.git
							</a>
						</td>
						<td>
							<span class="badge badge-secondary">b39bb312952c544826333b7578613c1e4d05c45c</span>
						</td>
						<td>
							<a class="btn btn-success">
								<i class="fa fa-check"></i>
								Enabled
							</a>
						</td>
						<td>
							<div class="btn-group">
								<a class="btn btn-primary">
									<i class="fa fa-edit"></i>
								</a>
								<a class="btn btn-danger">
									<i class="fa fa-times"></i>
								</a>
							</div>
						</td>
					</tr>
					<tr>
						<td>
							<span class="badge badge-secondary">Mod</span>
						</td>
						<td>3d_armor</td>
						<td>
							<span class="badge badge-primary">
								ContentDB
							</span>
						</td>
						<td>
							<a href="https://content.minetest.net/packages/stu/3d_armor/">
								https://content.minetest.net/packages/stu/3d_armor/
							</a>
						</td>
						<td>
							<span class="badge badge-secondary">
								2021-02-08
								<i class="fa fa-bell" style="color: yellow;"></i>
							</span>
							<a class="btn btn-success">
								<i class="fa fa-download"></i>
							</a>
						</td>
						<td>
							<a class="btn btn-secondary">
								<i class="fa fa-times"></i>
								Disabled
							</a>
						</td>
						<td>
							<div class="btn-group">
								<a class="btn btn-primary">
									<i class="fa fa-edit"></i>
								</a>
								<a class="btn btn-danger">
									<i class="fa fa-times"></i>
								</a>
							</div>
						</td>
					</tr>
				</tbody>
				<tbody>
					<tr>
						<td>
							<select class="form-control">
								<option>Mod</option>
								<option>Game</option>
								<option>Textures</option>
							</select>
						</td>
						<td></td>
						<td>
							<select class="form-control">
								<option>Git</option>
								<option>ContentDB</option>
							</select>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Source url"/>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Version"/>
						</td>
						<td>
						</td>
						<td>
							<a class="btn btn-success">
								<i class="fa fa-save"></i>
							</a>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};
