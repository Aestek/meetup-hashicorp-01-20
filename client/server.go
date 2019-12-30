package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Server struct {
	listen string
	rogue  bool

	stats        chan Stats
	currentStats Stats
}

func NewServer(listen string, stats chan Stats, rogue bool) *Server {
	return &Server{
		listen: listen,
		stats:  stats,
		rogue:  rogue,
	}
}

func (s *Server) Run() error {
	go func() {
		for stats := range s.stats {
			s.currentStats = stats
		}
	}()

	tmpl, err := template.New("index").Parse(indexTmpl)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, map[string]interface{}{
			"rogue": s.rogue,
		})
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.currentStats)
	})

	return http.ListenAndServe(s.listen, nil)
}

const indexTmpl = `
<html>

<head>
    <style>
        html, body {
			font-family: Arial, Helvetica, sans-serif;

			{{ if .rogue }}
			background: #fc4b4b;
			color: white;
			{{ end }}
        }

        h1 {
            max-width: 600px;
            margin: 50px auto;
            padding: 30px;
            border-bottom: 1px dashed black;
        }

        #stats {
            max-width: 600px;
            margin: 50px auto;
            padding: 30px;
        }


        #success,
        #errors,
        #body {
            max-width: 600px;
            min-height: 20px;
            background: rgba(0, 0, 0, 0.3);
            margin: 50px auto;
            padding: 30px;
            text-align: right;
            position: relative;
        }

        #success::before {
            content: 'Success:';
            display: block;
            float: left;

        }

        #errors::before {
            content: 'Errors:';
            display: block;
            float: left;
        }

        .delta {
            position: absolute;
            left: calc(100% + 10px);
            display: none;
            font-weight: bold;
            opacity: .8;
        }

        #body {
            text-align: justify;
            word-break: break-all;
            font-family: monospace, monospace;
        }
    </style>
    <script src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
    <script>
        $(function () {
            var lastSuccess = 0;
            var lastErrors = 0;
            function fetch() {
                $.ajax({
                    url: '/stats',
                    success: function (res) {
                        var data = JSON.parse(res);
                        var success = data.success || 0;
                        var errors = data.errors || 0;

                        $('#success .count').text(success);
                        $('#errors .count').text(errors);

                        if (success - lastSuccess) {
                            $('#success .delta').text('+' + (success - lastSuccess)).show().fadeOut();
                        }
                        if (errors - lastErrors) {
                            $('#errors .delta').text('+' + (errors - lastErrors)).show().fadeOut();
                        }
                        $('#body').text(data.last_body);

                        lastSuccess = success;
                        lastErrors = errors;
                    },
                });
            }
            fetch();
            setInterval(fetch, 1000);
        });
    </script>

    <title>{{ if .rogue }}Rogue client 😈{{ else }}Connect demo client{{ end }}</title>
</head>

<body>
    <h1>{{ if .rogue }}Rogue client 😈{{ else }}Connect demo client{{ end }}</h1>
    <div id="success">
        <span class="count"></span>
        <span class="delta"></span>
    </div>
    <div id="errors">
        <span class="count"></span>
        <span class="delta"></span>
    </div>
    <div id="body"></div>
</body>

</html>
`
