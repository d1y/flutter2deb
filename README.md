# flutter desktop to linux `deb`

```bash
go get github.com/d1y/flutter2deb

flutter2deb <abs path>
```

the package create build `.deb` template file

```log
debian
├── DEBIAN
│   └── control
└── usr
    ├── bin
    │   ├── data
    │   │   ├── flutter_assets
    │   │   └── icudtl.dat
    │   ├── lib
    │   │   ├── libapp.so
    │   │   ├── libflutter_linux_gtk.so
    │   │   └── liburl_launcher_linux_plugin.so
    │   └── spark_store
    └── share
        ├── applications
        │   └── spark_store.desktop
        └── icons
            └── spark_store.png
```

Read `pubspec.yaml` file get `name` | `description` | `version` Environment variable

And use `git config` get `username` and `email` Environment variable

You root need to have { `app.png` | `app.jpg` | `app.svg` } icon file(after add `*.desktop` file)

After this I created a `build.sh`, just run :)