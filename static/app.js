const appDiv = document.getElementById('app');

const tmpl = `
<div class="row">
    <div class="col-md-5">
        <div class="card">
            <div class="card-body">
                <h5 class="card-title">Total requests</h5>
                <div class="card-text">
                    <h3>\{{total_requests}}</h3>
                </div>
            </div>
        </div>
    </div>
    <div class="col-md-5">
        <div class="card">
            <div class="card-body">
                <h5 class="card-title">Average response time</h5>
                <div class="card-text">
                    <h3>\{{ average_response_time }} seconds</h3>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="row">
    <div class="col-md-5">
        <div class="card">
            <div class="card-body">
                <h5 class="card-title">Busiest days of the week</h5>
                <div class="card-text" style="width: 18rem">
                    <ul class="list-group list-group-flush">
                        {{#each requests_per_day}}
                        <li class="list-group-item">
                            \{{ this.id }} (\{{ this.number_of_requests }} requests)
                        </li>
                        {{/each }}
                    </ul>
                </div>
            </div>
        </div>
    </div>
    <div class="col-md-5">
        <div class="card">
            <div class="card-body">
                <h5 class="card-title">Busiest hours of day</h5>
                <div class="card-text" style="width: 18rem;">
                    <ul class="list-group list-group-flush">
                        {{#each requests_per_hour}}
                        <li class="list-group-item">
                            \{{ this.id }} (\{{ this.number_of_requests }} requests)
                        </li>
                        {{/each}}
                    </ul>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="row">
    <div class="col-md-5">
        <div class="card">
            <div class="card-body">
                <h5 class="card-title">Most visited routes</h5>
                <div class="card-text" style="width: 18rem;">
                    <ul class="list-group list-group-flush">
                        {{#each stats_per_route}}
                        <li class="list-group-item">
                            \{{ this.id.method }} \{{ this.id.url }} (\{{ this.number_of_requests }} requests)
                        </li>
                        {{/each}}
                    </ul>
                </div>
            </div>
        </div>
    </div>
</div>
`;

const template = Handlebars.compile(tmpl);

writeData = data => {
  appDiv.innerHTML = template(data);
};

axios
  .get('http://localhost:4000/api/analytics', {})
  .then(res => {
    console.log(res.data);
    writeData(res.data);
  })
  .catch(err => {
    console.error(err);
  });

const APP_KEY = 'PUSHER_APP_KEY';
const APP_CLUSTER = 'PUSHER_CLUSTER';

const pusher = new Pusher(APP_KEY, {
  cluster: APP_CLUSTER,
});

const channel = pusher.subscribe('analytics-dashboard');

channel.bind('data', data => {
  writeData(data);
});
