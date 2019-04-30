//Modify the backend url here
var API_URL = 'http://localhost:8000/'

Vue.mixin({
    data: function () {
        return {
            API_URL: API_URL
        }
    }
});

var result_parent = Vue.component('dropdown-url', {
    props: {
        urls: Array
    },
    template: `
    <div class="btn-group">
        <button class="btn btn-light btn-sm dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
            <slot></slot>
        </button>
        <div class="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <a v-for="url in urls" class="dropdown-item" _target="blank" v-bind:href="url">{{url}}</a>
        </div>
        </div>
    `
})

var result_keyword = Vue.component('result-keyword', {
    props: {
        keywords: Array
    },
    template: `
        <div>
        <span v-for="keyword in keywords" class="badge badge-light mr-2">{{keyword.word + " " + keyword.frequency}}</span>
        </div>
    `
})


var searchResult = Vue.component('search-result', {
    props: {
        title: String,
        url: String,
        date: String,
        score: Number

    },
    template: `
        <div class="card result">
            <div class="card-body">
                <p><a v-bind:href="url" class="title">{{title}}</a><span class="score">({{score.toFixed(3)}})</span></p>
                <p><a v-bind:href="url" class="url">{{url}}</a></p>
                <p class="date">{{date}}</p>
                <slot></slot>
            </div>
        </div>
    `
})

new Vue({
    el: '#app',
    data: {
        page: 1,
        errors: "",
        query: '',
        results: [],
        keywords: [],
        notfound: false,
        newQuery: []
    },
    methods: {
        checkQuery: function (e) {

            if (!this.query) {
                this.errors = "Search query cannot be empty";
            } else {
                this.errors = '';
                this.search(this.query)
            }

            if (!this.errors) {
                return true;
            }

            e.preventDefault();
        },

        search: function (query) {
            this.notfound = false;
            this.page = 1;
            this.results = [];
            axios.get(this.API_URL + "query/" + query)
                .then(res => {

                    this.results = res.data.documents;
                    if (!res.data.documents) {
                        this.notfound = true;
                    }
                    if(res.data.documents.length==0){
                        this.notfound = true;
                    }
                })
        },

        gotoSearch: function () {
            this.page = 1;
        },

        gotoKeywords: function () {
            this.newQuery = [];
            if (!this.keywords.length) {
                axios.get(this.API_URL + "wordList")
                    .then(res => {
                        this.keywords = res.data.words;
                    })
            }
            this.page = 2;
        },
        appendQuery: function (word) {
            if (!this.newQuery.includes(word))
                this.newQuery.push(word);
        },
        removeQuery: function (word) {
            var index = this.newQuery.indexOf(word);
            if (index > -1) {
                this.newQuery.splice(index, 1);
            }
        },
        selected: function (word) {
            if (this.newQuery.includes(word))
                return "btn-secondary disabled";
            else
                return "btn-light";
        },
        searchNewQuery: function () {

            this.query = this.newQuery.join(' ');
            this.search(this.query);

        },
        loadGraph: function (id) {
            d3.selectAll('#graph svg').remove();
            d3.json(API_URL+"graph/"+id, function (error, links) {
                
                var nodes = {};
                links = links.EdgesString
                // Compute the distinct nodes from the links.
                links.forEach(function (link) {
                    link.source = nodes[link.source] ||
                        (nodes[link.source] = { name: link.source });
                    link.target = nodes[link.target] ||
                        (nodes[link.target] = { name: link.target });
                    link.value = 1;
                });

                var width = 1100,
                    height =800;

                var force = d3.layout.force()
                    .nodes(d3.values(nodes))
                    .links(links)
                    .size([width, height])
                    .linkDistance(150)
                    .charge(-400)
                    .on("tick", tick)
                    .start();

                var svg = d3.select("#graph").append("svg")
                    .attr("width", width)
                    .attr("height", height);

                // build the arrow.
                svg.append("svg:defs").selectAll("marker")
                    .data(["end"])      // Different link/path types can be defined here
                    .enter().append("svg:marker")    // This section adds in the arrows
                    .attr("id", String)
                    .attr("viewBox", "0 -5 10 10")
                    .attr("refX", 15)
                    .attr("refY", -1.5)
                    .attr("markerWidth", 6)
                    .attr("markerHeight", 6)
                    .attr("orient", "auto")
                    .append("svg:path")
                    .attr("d", "M0,-5L10,0L0,5");

                // add the links and the arrows
                var path = svg.append("svg:g").selectAll("path")
                    .data(force.links())
                    .enter().append("svg:path")
                    //    .attr("class", function(d) { return "link " + d.type; })
                    .attr("class", "link")
                    .attr("marker-end", "url(#end)");

                // define the nodes
                var node = svg.selectAll(".node")
                    .data(force.nodes())
                    .enter().append("g")
                    .attr("class", "node")
                    .call(force.drag);

                // add the nodes
                node.append("circle")
                    .attr("r", 5);

                // add the text 
                node.append("text")
                    .attr("x", 12)
                    .attr("dy", ".35em")
                    .text(function (d) { return d.name; });

                // add the curvy lines
                function tick() {
                    path.attr("d", function (d) {
                        var dx = d.target.x - d.source.x,
                            dy = d.target.y - d.source.y,
                            dr = Math.sqrt(dx * dx + dy * dy);
                        return "M" +
                            d.source.x + "," +
                            d.source.y + "A" +
                            dr + "," + dr + " 0 0,1 " +
                            d.target.x + "," +
                            d.target.y;
                    });

                    node
                        .attr("transform", function (d) {
                            return "translate(" + d.x + "," + d.y + ")";
                        });
                }
            });
        }

    } 
})
