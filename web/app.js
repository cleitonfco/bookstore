/** @jsx React.DOM */
 
var App = React.createClass({
  getInitialState: function() {
    return {
      total: 0, 
      query: "", 
      books: [], 
      page: 1
    };
  },
  handlePage: function(p) {
    var numPages = Math.ceil(this.state.total / this.props.pageSize)
    if (p > 0 && p <= numPages) {
      this.loadSearch(this.state.query, p);
    }
  }, 
  handleSearch: function(q) {
    if (q != this.state.query) {
      this.loadSearch(q, 1);
    }
  }, 
  loadSearch: function(query, page) {
    $.ajax({
      url: '/books',
      data: { 'q': query, 'p': page }, 
      dataType: 'json',
      success: function(data) {
        this.setState({
          total: data.total, 
          books: data.books, 
          query: query, 
          page: page
        });
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/books', status, err.toString());
      }.bind(this)
    });
  },
  componentDidMount: function() {
    this.loadSearch(this.state.query, this.state.page);
  }, 
  render: function() {
    return <div>
      <FormSearch query={this.state.query} submitSearch={this.handleSearch} />
      {this.state.total > 0 ?
        <div>
          <BookListing books={this.state.books} />
          <Pagination total={this.state.total} pageSize={this.props.pageSize} currentPage={this.state.page} setPage={this.handlePage} />
        </div> : 
        <p>Não foi possível localizar nenhum produto para o termo indicado</p>
      }
    </div>
  }
})
 
var FormSearch = React.createClass({
  getInitialState: function() {
    return {
      term: this.props.query
    };
  },
  onChange: function(e) {
    this.setState({
      term: e.target.value
    });
  },
  handleSubmit: function(e) {
    e.preventDefault();
    var term = this.refs.term.getDOMNode().value.trim();
    this.setState({
      term: term
    });
    this.props.submitSearch(term);
  },
  render: function() {
    return <div className="row">
      <div className="col-md-10 col-md-offset-1 col-lg-8 col-lg-offset-2">
        <form className="form-group" onSubmit={this.handleSubmit}>
          <div className="input-group input-group-lg">
            <input onChange={this.onChange} type="search" name="q" className="form-control" ref="term" value={this.state.term} />
            <span className="input-group-btn">
              <button className="btn btn-primary" type="button" onClick={this.handleSubmit}>Buscar</button>
            </span>
          </div>
        </form>
      </div>
    </div>
  }
})
 
var PaginationItem = React.createClass({
  setActive: function () {
    this.props.setPage(this.props.page);
  },
  render: function() {
    return <li className={this.props.active ? "active" : ""} onClick={this.setActive}>
      <a>{this.props.page}</a>
    </li>
  }
});

var Pagination = React.createClass({
  handleNext: function() {
    this.props.setPage(this.props.currentPage + 1);
  },
  handlePrevious: function() {
    this.props.setPage(this.props.currentPage - 1);
  },

  render: function() {
    var numPages = Math.ceil(this.props.total / this.props.pageSize)
      , item     = 1
      , pages    = []
    while (item <= numPages) {
      pages.push(item);
      item++;
    }
    var paginationItems = pages.map(function (p) {
      return <PaginationItem page={p} active={p==this.props.currentPage} setPage={this.props.setPage} key={p} />;
    }.bind(this));

    var previousLink = (this.props.currentPage > 1) ? 
      <li><a onClick={this.handlePrevious}>Previous</a></li> :
      <li className="disabled"><span>Previous</span></li>;
    var nextLink = (this.props.currentPage < numPages) ?
      <li><a onClick={this.handleNext}>Next</a></li> : 
      <li className="disabled"><span>Next</span></li>;

    return <ul className="pagination">
      {previousLink}
      {paginationItems}
      {nextLink}
    </ul>
  }
});

var BookListing = React.createClass({
  render: function() {
    var books = this.props.books.map(function(book) {
      return <li key={book.id} className="col-sm-4 col-md-3 col-lg-2">
        <img src={book.image} alt={book.title} className="image-cover" />
        <p className="book-name">{book.title}</p>
        <p className="author-name">{book.author}</p>
        <a className="btn btn-primary">Comprar</a>
      </li>
    })
    return <div className="row">
      <div>
        <h2>Livros</h2>
        <ul className="books">
          {books}
        </ul>
      </div>
    </div>
  }
})

React.renderComponent(<App pageSize={4} />, document.getElementById('componentPage'))