axios
  .get('http://localhost:4000/api/analytics', {})
  .then(res => {
    console.log(res);
  })
  .catch(err => {
    console.error(err);
  });

const app = new Vue({
  el: '#app',
});

Vue.component('analytics', {
  data: window.analytics,
  template: `

    <div class="row">
        <div class="col-md-5">
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">Total requests</h5>
                    <div class="card-text">
                        <h3>{{ totalRequests }}</h3>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-md-5">
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">Average response time</h5>
                    <div class="card-text">
                        <h3 id="averageResponse">{{ averageResponseTime }} seconds</h3>
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
                    <div class="card-text" style="width: 18rem;" v-for="day in busiestDays">
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item">
                                {{ day._id }} ({{ day.numberOfRequests }} requests)
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-md-5">
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">Busiest hours of day</h5>
                    <div class="card-text" style="width: 18rem;" v-for="hour in busiestHours">
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item">
                                {{ hour._id }} ({{ hour.numberOfRequests }} requests)
                            </li>
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
                    <div class="card-text" style="width: 18rem;" v-for="route in statsPerRoute">
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item">
                                {{ route._id.method }} {{ route._id.url }} ({{ route.numberOfRequests }} requests)
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-md-5">
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">Slowest routes</h5>
                    <div class="card-text" style="width: 18rem;" v-for="route in statsPerRoute">
                        <ul class="list-group list-group-flush">
                            {{ route._id.method }} {{ route._id.url }} ({{ route.responseTime || 0}} s)
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>


	    `,
});
