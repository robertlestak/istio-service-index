var app = new Vue({
  el: '#app',
  data: {
    categories: ["media", "tools"],
    services: [],
    searchInput: "",
    loading: true,
  },
  computed: {
    serviceList() {
      return []
    }
  },
  methods: {
    async getServices() {
      try {
        this.loading = true
        let res = await axios.get("http://i.int/api")
        this.loading = false
        this.services = res.data
        this.getCategories()
      } catch (e) {
        console.log(e)
      }
    },
    goService(s) {
      location.href = 'http://'+s.hosts[0]
    },
    getCategories() {
      for (i in this.services) {
        if (this.categories.indexOf(this.services[i].category) == -1) {
          this.categories.push(this.services[i].category)
        }
      }
    },
    categoryIcon(c) {
      let b = "fa fa-3x "
      switch (c) {
        case "communication":
          return b += "fa-paper-plane"
          break
        case "media":
          return b += "fa-music"
          break
        case "tools":
          return b += "fa-wrench"
          break
        case "home":
          return b += "fa-home"
          break
        case "finance":
          return b += "fa-landmark"
          break
        default:
          return b += "fa-"+c
      }
      return b
    },
    serviceContainsString(svc, s) {
      if (s === "") return true
      s = s.toLowerCase()
      if (svc.category.toLowerCase().includes(s)) {
        return true
      } else if (svc.name.toLowerCase().includes(s)) {
        return true
      } else if (svc.description.toLowerCase().includes(s)) {
        return true
      }
    },
    categoryServices(c) {
      let cs = []
      for (i in this.services) {
        if (this.services[i].category === c && this.serviceContainsString(this.services[i], this.searchInput)) {
          cs.push(this.services[i])
        }
      }
      return cs
    }
  },
  mounted() {
    this.getServices()
  }
})
