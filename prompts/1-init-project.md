I want to build a app that can help me:

1. Search moviews and TV shows by keyword or imdb id, with selected prowlarr indexes and catagories. have button to download the resource via move the torrent to designated dir.
2. Watch search queries (wishlist) in prowlarr, maybe search every hour, if this is found auto download via save the torrent to designed dir (can be configable).
3. Watch the download complete dir (configable), has a page list completed downloads, with a text box for rename downloaded files.

This should include:

1. web frontend writted in React and Tailwind, you can place this in web/, and you can assume I use tool `pnpm` instead of `npm`.
2. backend writted in Golang, you can place this in backend/
3. backend and frontend use openapi to communicate, put define in apis/

Backend:

2. Use sqlite to save wishlist, downloads
3. Backend should support google login user token as access token
4. Backend should support config:
    - list of downloader dirs: name - dir mapping
    - prorarr api url, and api key
    - google project id for google login, list of allowed e-mail addresses
    - sqlite db path
    - api key of tmdb
    - download complete dir
    - list of libary dirs: name - dir mapping
5. when backend start download a torront from wishlist, it should marked it as "done" wishlist. and download the torront file from prowlarr to selected downloader dir. and add it to downloads table.
6. backend should watch the download complete dir, check if the download is completed, it maybe new file in the completed dir or in a sub dir. if it is completed, it should marked it as "downloaded" in downloads.

Frontend should:

1. only support communicated with backend via openapi, and display data from backend.
2. frontend should use google login user token as access token to authenticate with backend. google project id should be configable.

Frontend ui should list of tabs on the left pane, and main content on the right:

1. Search

   - dropdown to select use keyword or imdb id to search
   - text input box for keyword or imdb id
   - drowpdown to select prowlarr indexes, can be multi select, frontend can get list of prowlarr indexes should via api to backend, frontend should cache this list in browser local storage, a refreash button to refreash this list, and auto fetch this list if not eixst in local storage
   - drowpdown to select catagories , can be multi select, frontend can get list of catagories should via api to backend, frontend should cache this list in browser local storage, a refreash button to refreash this list, and auto fetch this list if not eixst in local storage
   - a text box for size, unit is GB, allow user to input a float number, backend should only return download size larger than this number
   - button to search, call api to backend
   - a button add the current search to wishlist
   - list search results, each item should include:
        1. title
        2. indexer
        3. size
        4. download button, click on download button should popup a dialog to select which downloader to use, then call api to backend to download this item

2. Wishlist

   - list current wishlist, each item should include the search queires and status of this item. status can be: `active`, `done`.
   - a button to remove item from wishlist
   - a button to edit this item from wishlist, so user can edit the search queries, for example we have done for "TV Show S01E01", user can edit this to "TV Show S01E02"

3. Downloaded

   - list downloaded torrent, each item should be expandable, and show the file list of this download, backend can use get this info via parse the torent file.
   - a text box for copy each downloaded files to a libary (dropdown to select), then user can input a new name for this file. a button to submit all of them to backend. backend should copy the downloaded files to libary dir with new name. new name can be "dir/file", if dir not exist, backend should create it.
   - after copy submitted, backend should marked this item as "copied" in downloads.

-------------------------

The prompt I use above is basically useless as it does not provide clear instruction to generate code. Let's ask AI to break down the requirements and generate code for each part.

-------------------------

Here is a bad product requirement doc, please me rewrite, please descript all restful endpoint need to be used between backend and frontend, how backend impl each api, and how frontend use them and present them. You can list all restful endpoint first so I can check if any of them are missing.

-------------------------

Please ignore Authentication and Config, gen openapi file
