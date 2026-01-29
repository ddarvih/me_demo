#ifndef Demo_TensorLikeObject_H
#define Demo_TensorLikeObject_H
#include <array>
#include <limits>
#include <optional>
#include <stdexcept>

template< typename T, size_t N >
class Demo_TensorLikeObject;

template< typename T, size_t N, typename C, size_t F >
class Demo_TensorLikeObjectViewBase
{
  public:
	using arrN = std::array< size_t, N >;
	using iterator_category = std::contiguous_iterator_tag;
	using difference_type = std::ptrdiff_t;
	using value_type = C;
	using pointer = C*;
	using reference = C&;
	using size_type = std::size_t;

	friend class Demo_TensorLikeObject< T, N >;

	Demo_TensorLikeObjectViewBase() noexcept = default;

	Demo_TensorLikeObjectViewBase(pointer data, const arrN* dim_sz, const arrN* dim_steps, size_t base_offset, size_t cur_offset) :
		m_data(data), m_dim_sz(dim_sz), m_dim_steps(dim_steps), m_dims_viewed_fix(F), m_base_offset(base_offset),
		m_cur_offset(cur_offset)
	{
	}

	Demo_TensorLikeObjectViewBase(pointer data, const arrN* dim_sz, const arrN* dim_steps, const std::array< size_t, F >& viewed_indexes) :
		m_data(data), m_dim_sz(dim_sz), m_dim_steps(dim_steps), m_dims_viewed_fix(F), m_base_offset(0), m_cur_offset(0)
	{
		for (size_t i = 0; i < F; i++)
			m_base_offset += viewed_indexes[i] * (*m_dim_steps)[i];
	}

	template< size_t Fcur = F >
	constexpr Demo_TensorLikeObjectViewBase(pointer data, const arrN* dim_sz, const arrN* dim_steps, const std::array< size_t, Fcur >& viewed_indexes)
		requires(Fcur + 1 == F && Fcur < N)
		:
		m_data(data), m_dim_sz(dim_sz), m_dim_steps(dim_steps), m_dims_viewed_fix(F), m_base_offset(0), m_cur_offset(0)
	{
		for (size_t i = 0; i < Fcur; i++)
			m_base_offset += viewed_indexes[i] * (*m_dim_steps)[i];
	}

	Demo_TensorLikeObjectViewBase(const Demo_TensorLikeObjectViewBase& other) :
		m_data(other.m_data), m_dim_sz(other.m_dim_sz), m_dim_steps(other.m_dim_steps), m_dims_viewed_fix(other.m_dims_viewed_fix),
		m_base_offset(other.m_base_offset), m_cur_offset(other.m_cur_offset), m_arr_carry(other.m_arr_carry)
	{
		if (m_arr_carry)
		{
			m_dim_sz = &m_arr_carry->first;
			m_dim_steps = &m_arr_carry->second;
		}
	}

	Demo_TensorLikeObjectViewBase(Demo_TensorLikeObjectViewBase&& other) noexcept :
		m_data(other.m_data), m_dim_sz(other.m_dim_sz), m_dim_steps(other.m_dim_steps), m_dims_viewed_fix(other.m_dims_viewed_fix),
		m_base_offset(other.m_base_offset), m_cur_offset(other.m_cur_offset), m_arr_carry(std::move(other.m_arr_carry))
	{
		if (m_arr_carry)
		{
			m_dim_sz = &m_arr_carry->first;
			m_dim_steps = &m_arr_carry->second;
		}
	}

	Demo_TensorLikeObjectViewBase& operator=(const Demo_TensorLikeObjectViewBase& other) noexcept
	{
		if (this == other)
			return *this;
		m_data = other.m_data;
		m_dims_viewed_fix = other.m_dims_viewed_fix;
		m_base_offset = other.m_base_offset;
		m_cur_offset = other.m_cur_offset;
		m_arr_carry = other.m_arr_carry;
		if (m_arr_carry)
		{
			m_dim_sz = &m_arr_carry->first;
			m_dim_steps = &m_arr_carry->second;
		}
		return *this;
	}

	Demo_TensorLikeObjectViewBase& operator=(Demo_TensorLikeObjectViewBase&& other) noexcept
	{
		if (this == other)
			return *this;
		m_arr_carry.reset();
		m_data = other.m_data;
		m_dims_viewed_fix = other.m_dims_viewed_fix;
		m_arr_carry = std::move(other.m_arr_carry);
		if (m_arr_carry)
		{
			m_dim_sz = &m_arr_carry->first;
			m_dim_steps = &m_arr_carry->second;
		}
		else
		{
			m_dim_sz = other.m_dim_sz;
			m_dim_steps = other.m_dim_steps;
		}
		return *this;
	}

	~Demo_TensorLikeObjectViewBase() noexcept = default;

	explicit Demo_TensorLikeObjectViewBase(const Demo_TensorLikeObject< T, N >& nda) :
		m_data(nda.m_data), m_dim_sz(&nda.m_dim_sz), m_dim_steps(&nda.m_dim_steps), m_dims_viewed_fix(0),
		m_base_offset(0), m_cur_offset(0)
	{
	}

	template< typename SizeList >
	Demo_TensorLikeObjectViewBase(value_type* data, const SizeList& sizes) :
		m_data(data), m_dim_sz{}, m_dim_steps(nullptr), m_dims_viewed_fix(0), m_base_offset(0), m_cur_offset(0)
	{
		m_arr_carry = std::pair< arrN, arrN >{};
		m_dim_sz = m_arr_carry->first;
		m_dim_steps = m_arr_carry->second;
		for (size_t i = 0; i < sizes.size(); i++)
			(*m_dim_sz)[i] = sizes[i];
		size_t cur_step = 1;
		if (N > 0)
		{
			for (int i = N - 1; i >= 0; i--)
			{
				(*m_dim_steps)[i] = cur_step;
				cur_step *= (*m_dim_sz)[i];
			}
		}
	}

	constexpr Demo_TensorLikeObjectViewBase(pointer data, const arrN* dim_sz, const arrN* dim_steps, size_t viewed_index = 0)
		requires(F == 0)
		: m_data(data), m_dim_sz(dim_sz), m_dim_steps(dim_steps), m_dims_viewed_fix(0), m_base_offset(0)
	{
		m_cur_offset = viewed_index * (*dim_steps)[0];
	}

	template< size_t Fcur = F >
	Demo_TensorLikeObjectViewBase(pointer data, const arrN* dim_sz, const arrN* dim_steps, size_t dims_viewed_fix, const std::array< size_t, Fcur >& viewed_indexes) :
		m_data(data), m_dim_sz(dim_sz), m_dim_steps(dim_steps), m_dims_viewed_fix(dims_viewed_fix), m_base_offset(0),
		m_cur_offset(0)
	{
		for (size_t i = 0; i < Fcur; i++)
			m_base_offset += viewed_indexes[i] * (*m_dim_steps)[i];
	}

	template< std::input_iterator It >
	Demo_TensorLikeObjectViewBase(It it1, It it2)
		requires std::input_iterator< It >
		: m_dims_viewed_fix(0), m_base_offset(0), m_cur_offset(0)
	{
		m_arr_carry = std::pair< arrN, arrN >{};
		m_dim_sz = &m_arr_carry->first;
		m_dim_steps = &m_arr_carry->second;
		if (it1 == it2)
		{
			m_data = nullptr;
			return;
		}
		m_data = std::addressof(*it1);
		if constexpr (N > 0)
		{
			m_dim_sz->fill(0);
			m_dim_steps->fill(0);
			(*m_dim_sz)[0] = std::distance(it1, it2);
			(*m_dim_steps)[0] = 1;
		}
		if constexpr (std::is_same_v< std::remove_cv_t< It >, Demo_TensorLikeObjectViewBase< T, N, T, 0 > > ||
					  std::is_same_v< std::remove_cv_t< It >, Demo_TensorLikeObjectViewBase< T, N, const T, 0 > >)
		{
			m_dim_sz = it1.m_dim_sz;
			m_dim_steps = it1.m_dim_steps;
			m_base_offset = it1.m_base_offset;
			m_cur_offset = it1.m_cur_offset;
		}
	}

	Demo_TensorLikeObjectViewBase& operator++() noexcept
	{
		if (m_dims_viewed_fix < N)
			m_cur_offset += (*m_dim_steps)[m_dims_viewed_fix];
		return *this;
	}
	Demo_TensorLikeObjectViewBase& operator--() noexcept
	{
		if (m_dims_viewed_fix < N)
			m_cur_offset -= (*m_dim_steps)[m_dims_viewed_fix];
		return *this;
	}
	Demo_TensorLikeObjectViewBase& operator++(int) noexcept
	{
		Demo_TensorLikeObjectViewBase cur = *this;
		++(*this);
		return cur;
	}
	Demo_TensorLikeObjectViewBase& operator--(int) noexcept
	{
		Demo_TensorLikeObjectViewBase cur = *this;
		--(*this);
		return cur;
	}

	Demo_TensorLikeObjectViewBase& operator+=(difference_type n) noexcept
	{
		if (m_dims_viewed_fix < N)
		{
			difference_type cur_ind = static_cast< difference_type >(m_cur_offset) + n;
			m_cur_offset = static_cast< size_t >(cur_ind);
		}
		return *this;
	}
	Demo_TensorLikeObjectViewBase& operator-=(difference_type n) noexcept
	{
		*this += -n;
		return *this;
	}
	Demo_TensorLikeObjectViewBase operator+(difference_type n) const noexcept
	{
		Demo_TensorLikeObjectViewBase cur = *this;
		cur += n;
		return cur;
	}
	Demo_TensorLikeObjectViewBase operator-(difference_type n) const noexcept
	{
		Demo_TensorLikeObjectViewBase cur = *this;
		cur -= n;
		return cur;
	}
	bool operator==(const Demo_TensorLikeObjectViewBase& other) const noexcept
	{
		return m_data == other.m_data && m_dim_sz == other.m_dim_sz && m_dim_steps == other.m_dim_steps &&
			   m_dims_viewed_fix == other.m_dims_viewed_fix && m_base_offset == other.m_base_offset &&
			   m_cur_offset == other.m_cur_offset;
	}
	bool operator!=(const Demo_TensorLikeObjectViewBase& other) const noexcept { return !(*this == other); }

	difference_type operator-(const Demo_TensorLikeObjectViewBase& other) const noexcept
	{
		return static_cast< difference_type >(m_cur_offset) - static_cast< difference_type >(other.m_cur_offset);
	}

	auto operator<=>(const Demo_TensorLikeObjectViewBase& other) const noexcept
	{
		return (m_base_offset + m_cur_offset) <=> (other.m_base_offset + other.m_cur_offset);
	}

	template< typename SizeList >
	auto at(const SizeList& ind_list) const
	{
		if (ind_list.size() != N)
			throw std::out_of_range("Demo_TensorLikeObjectViewBase::at !=N");
		size_t cur_ind = 0;
		for (size_t d = 0; d < N; ++d)
		{
			if (ind_list[d] >= (*m_dim_sz)[d])
				throw std::out_of_range("Demo_TensorLikeObjectViewBase::at");
			cur_ind += ind_list[d] * (*m_dim_steps)[d];
		}
		return m_data[cur_ind];
	}

	auto begin() const noexcept { return Demo_TensorLikeObjectViewBase(m_data, m_dim_sz, m_dim_steps, m_base_offset, 0); }
	auto end() const noexcept { return begin() + ((*m_dim_sz)[m_dims_viewed_fix] * (*m_dim_steps)[m_dims_viewed_fix]); }

	auto cbegin() const noexcept { return Demo_TensorLikeObjectViewBase(m_data, m_dim_sz, m_dim_steps, m_base_offset, 0); }
	auto cend() const noexcept
	{
		return cbegin() + ((*m_dim_sz)[m_dims_viewed_fix] * (*m_dim_steps)[m_dims_viewed_fix]);
	}

	size_type count() const noexcept { return m_dims_viewed_fix < N ? (*m_dim_sz)[m_dims_viewed_fix] : 1; }
	size_type total_count() const noexcept
	{
		size_t total_counter = 1;
		for (size_t i = m_dims_viewed_fix; i < N; i++)
			total_counter *= (*m_dim_sz)[i];
		return total_counter;
	}
	size_type dim() const noexcept { return N - m_dims_viewed_fix; }
	bool is_equal(const Demo_TensorLikeObjectViewBase& other) const noexcept
	{
		if (m_dim_sz != other.m_dim_sz || m_dim_steps != other.m_dim_steps || m_dims_viewed_fix != other.m_dims_viewed_fix ||
			m_base_offset != other.m_base_offset || m_cur_offset != other.m_cur_offset)
			return false;
		return std::equal(begin(), end(), other.begin(), other.end());
	}

	constexpr auto
		operator[](difference_type ind) const -> std::conditional_t< F + 1 == N, reference, Demo_TensorLikeObjectViewBase< T, N, C, F + 1 > >
		requires(F < N)
	{
		if (ind < 0 || ind >= static_cast< difference_type >((*m_dim_sz)[m_dims_viewed_fix]))
			throw std::out_of_range("Demo_TensorLikeObjectViewBase::operator[]");
		size_t base_offset = m_base_offset + ind * (*m_dim_steps)[m_dims_viewed_fix];
		if constexpr (F + 1 == N)
			return m_data[base_offset];
		else
			return Demo_TensorLikeObjectViewBase< T, N, C, F + 1 >(m_data, m_dim_sz, m_dim_steps, base_offset, 0);
	}

	constexpr auto operator*() const -> std::conditional_t< F + 1 == N, reference, Demo_TensorLikeObjectViewBase< T, N, C, F + 1 > >
	{
		if constexpr (F + 1 == N)
		{
			if (m_cur_offset >= (*m_dim_sz)[F] * (*m_dim_steps)[F])
				throw std::out_of_range("operator*() - View");
			return m_data[m_base_offset + m_cur_offset];
		}
		else
		{
			return Demo_TensorLikeObjectViewBase< T, N, C, F + 1 >(m_data, m_dim_sz, m_dim_steps, m_base_offset + m_cur_offset, 0);
		}
	}

	pointer operator->() const
	{
		if constexpr (N != 1)
			throw std::out_of_range("Demo_TensorLikeObjectViewBase::operator->() N!=1");
		else
			return m_data + (m_cur_offset + m_base_offset);
	}

	operator reference() const
		requires(F == N)
	{
		return this;
	}

  protected:
	value_type* m_data;
	const arrN* m_dim_sz;
	const arrN* m_dim_steps;
	size_t m_dims_viewed_fix;
	size_t m_base_offset;
	size_t m_cur_offset;

	std::optional< std::pair< arrN, arrN > > m_arr_carry;
};

template< typename T, size_t N, size_t F >
using Demo_TensorLikeObjectView = Demo_TensorLikeObjectViewBase< T, N, T, F >;

template< typename T, size_t N, size_t F >
using Demo_TensorLikeObjectConstView = Demo_TensorLikeObjectViewBase< T, N, const T, F >;

namespace parse_inner_utils
{
	template< typename I, size_t F >
	struct parseInnerInitL
	{
		using type_n = std::initializer_list< typename parseInnerInitL< I, F - 1 >::type_n >;
	};
	template< typename I >
	struct parseInnerInitL< I, 1 >
	{
		using type_n = std::initializer_list< I >;
	};
	template< typename I >
	struct parseInnerInitL< I, 0 >
	{
		using type_n = I;
	};
}	 // namespace parse_inner_utils

template< typename T, size_t N >
class Demo_TensorLikeObject
{
  public:
	using arrN = std::array< size_t, N >;
	template< typename I >
	using initL = std::initializer_list< I >;

	using value_type = T;
	using reference = T&;
	using const_reference = const T&;
	using iterator = Demo_TensorLikeObjectView< T, N, 0 >;
	using const_iterator = Demo_TensorLikeObjectConstView< T, N, 0 >;
	using difference_type = std::ptrdiff_t;
	using size_type = size_t;

	using view = Demo_TensorLikeObjectView< T, N, 0 >;
	using const_view = Demo_TensorLikeObjectConstView< T, N, 0 >;

	explicit Demo_TensorLikeObject(const arrN& sizes) : m_dimensions(N), m_size(1)
	{
		pre_fill_inits(sizes);
		fill_data_with_inits(m_size, T{});
	}

	Demo_TensorLikeObject(const arrN& sizes, const T& val) : m_size(1), m_dimensions(N)
	{
		pre_fill_inits(sizes);
		fill_data_with_inits(m_size, val);
	}

	Demo_TensorLikeObject(const arrN& sizes, T* elements) : m_dimensions(N), m_size(1)
	{
		pre_fill_inits(sizes);
		fill_data_with_inits(m_size, elements);
	}

	Demo_TensorLikeObject(typename parse_inner_utils::parseInnerInitL< T, N >::type_n inits) : m_dimensions(N)
	{
		arrN sizes = {};
		count_sz< 0 >(sizes, inits);
		check_sz(sizes);
		m_data = static_cast< T* >(operator new(m_size * sizeof(T)));
		size_t cur_pos = 0;
		try
		{
			form_init_data< 0 >(m_data, inits, cur_pos);
			if (cur_pos != m_size)
				throw std::invalid_argument("Demo_TensorLikeObject constr initL - not fitting size");
		} catch (...)
		{
			act_erase(cur_pos);
			operator delete(m_data);
			throw;
		}
	}

	template< size_t Fcur >
		requires(!std::is_const_v< T >)
	Demo_TensorLikeObject(Demo_TensorLikeObjectView< T, N, Fcur > v_view) :
		m_size(v_view.total_count()), m_dimensions(N), m_dim_sz(*v_view.m_dim_sz), m_dim_steps(*v_view.m_dim_steps)
	{
		fill_data_with_inits(m_size, v_view.begin());
	}

	template< size_t Fcur >
	Demo_TensorLikeObject(Demo_TensorLikeObjectConstView< T, N, Fcur > c_view) :
		m_size(c_view.total_count()), m_dimensions(N), m_dim_sz(*c_view.m_dim_sz), m_dim_steps(*c_view.m_dim_steps)
	{
		fill_data_with_inits(m_size, c_view.begin());
	}

	template< std::input_iterator It >
	Demo_TensorLikeObject(It it1, It it2)
		requires std::input_iterator< It >
		: m_dimensions(N), m_size(0), m_data(nullptr)
	{
		m_size = std::distance(it1, it2);
		if (m_size == 0)
		{
			m_dim_sz = {};
			m_dim_steps = {};
			return;
		}
		if constexpr (std::is_same_v< std::remove_cv_t< It >, view > || std::is_same_v< std::remove_cv_t< It >, const_view >)
		{
			m_dim_sz = *(it1.m_dim_sz);
			m_dim_steps = *(it1.m_dim_steps);
		}
		else
		{
			m_dim_sz.fill(0);
			m_dim_steps.fill(0);
			m_dim_sz[0] = m_size;
			m_dim_steps[0] = 1;
		}
		fill_data_with_inits(m_size, it1);
	}

	void swap(Demo_TensorLikeObject& other) noexcept
	{
		if (this == &other)
			return;
		using std::swap;
		swap(m_data, other.m_data);
		swap(m_size, other.m_size);
		swap(m_dimensions, other.m_dimensions);
		swap(m_dim_sz, other.m_dim_sz);
		swap(m_dim_steps, other.m_dim_steps);
	}

	friend void swap(Demo_TensorLikeObject< T, N >& ndarr1, Demo_TensorLikeObject< T, N >& ndarr2) noexcept { ndarr1.swap(ndarr2); }

	Demo_TensorLikeObject() noexcept : m_data(nullptr), m_size(0), m_dimensions(N), m_dim_sz{}, m_dim_steps{} {}

	Demo_TensorLikeObject(const Demo_TensorLikeObject< T, N >& other) :
		m_data(nullptr), m_size(other.m_size), m_dimensions(other.m_dimensions), m_dim_sz(other.m_dim_sz),
		m_dim_steps(other.m_dim_steps)
	{
		fill_data_with_inits(m_size, other.m_data);
	}

	Demo_TensorLikeObject& operator=(const Demo_TensorLikeObject< T, N >& other) noexcept
	{
		if (this != &other)
		{
			Demo_TensorLikeObject< T, N > tmp(other);
			swap(*this, tmp);
		}
		return *this;
	}

	Demo_TensorLikeObject(Demo_TensorLikeObject< T, N >&& other) noexcept :
		m_data(other.m_data), m_size(other.m_size), m_dimensions(other.m_dimensions), m_dim_sz(other.m_dim_sz),
		m_dim_steps(other.m_dim_steps)
	{
		other.m_size = 0;
		other.m_data = nullptr;
	}

	Demo_TensorLikeObject& operator=(Demo_TensorLikeObject< T, N >&& other) noexcept
	{
		if (this != &other)
		{
			act_erase(m_size);
			operator delete(m_data);
			m_data = nullptr;
			m_size = 0;
			this->swap(other);
		}
		return *this;
	}

	~Demo_TensorLikeObject() noexcept
	{
		act_erase(m_size);
		operator delete(m_data);
	}

	bool is_equal(const Demo_TensorLikeObject& rhs) const noexcept
	{
		if (this != &rhs)
		{
			if (m_size != rhs.m_size || !std::equal(m_dim_sz.begin(), m_dim_sz.end(), rhs.m_dim_sz.begin()))
				return false;
			return std::equal(m_data, m_data + m_size, rhs.m_data);
		}
		return true;
	}

	T* data() noexcept { return m_data; }
	const T* data() const noexcept { return m_data; }

	size_type count() const noexcept { return m_dim_sz[0]; }
	size_type total_count() const noexcept { return m_size; }
	size_type dim() const noexcept { return m_dimensions; }

	constexpr auto operator[](size_type i) -> std::conditional_t< N == 1, reference, Demo_TensorLikeObjectView< T, N, 1 > >
	{
		if (i >= m_dim_sz[0])
			throw std::out_of_range("Demo_TensorLikeObject::operator[]");
		if constexpr (N == 1)
		{
			return m_data[i * m_dim_steps[0]];
		}
		else
		{
			std::array< size_t, 1 > cur_viewed_indexes = { i };
			return Demo_TensorLikeObjectView< T, N, 1 >(m_data, &m_dim_sz, &m_dim_steps, cur_viewed_indexes);
		}
	}

	constexpr auto operator[](size_type i) const -> std::conditional_t< N == 1, const_reference, Demo_TensorLikeObjectConstView< T, N, 1 > >
	{
		if (i >= m_dim_sz[0])
			throw std::out_of_range("Demo_TensorLikeObject::operator[]");
		if constexpr (N == 1)
		{
			return m_data[i * m_dim_steps[0]];
		}
		else
		{
			std::array< size_t, 1 > cur_viewed_indexes = { i };
			return Demo_TensorLikeObjectConstView< T, N, 1 >(m_data, &m_dim_sz, &m_dim_steps, cur_viewed_indexes);
		}
	}

	T& at(std::initializer_list< size_type > indexes) { return at_inner< T& >(indexes); }

	const T& at(std::initializer_list< size_type > indexes) const { return at_inner< const T& >(indexes); }

	view begin() noexcept { return Demo_TensorLikeObjectView< T, N, 0 >(m_data, &m_dim_sz, &m_dim_steps, 0); }

	view end() noexcept { return begin() + m_dim_sz[0] * m_dim_steps[0]; }

	const_view cbegin() const noexcept { return Demo_TensorLikeObjectConstView< T, N, 0 >(m_data, &m_dim_sz, &m_dim_steps, 0); }

	const_view cend() const noexcept { return cbegin() + m_dim_sz[0] * m_dim_steps[0]; }

	const_view begin() const noexcept { return Demo_TensorLikeObjectConstView< T, N, 0 >(m_data, &m_dim_sz, &m_dim_steps, 0); }

	const_view end() const noexcept
	{
		return Demo_TensorLikeObjectConstView< T, N, 0 >(m_data, &m_dim_sz, &m_dim_steps, m_dim_sz[0] * m_dim_steps[0]);
	}

	template< typename SizeList >
	Demo_TensorLikeObjectView< T, N, 0 > reshape(const SizeList& sizes) const
	{
		size_t check_counter = 1;
		for (const auto& sz : sizes)
			check_counter *= sz;
		if (check_counter != m_size)
			throw std::invalid_argument("Demo_TensorLikeObject::reshape - didn't pass the total amount check");
		return Demo_TensorLikeObjectView< T, N, 0 >(m_data, sizes);
	}

	bool operator==(const Demo_TensorLikeObject< T, N >& other) const { return is_equal(other); }
	bool operator!=(const Demo_TensorLikeObject< T, N >& other) const { return !is_equal(other); }

	size_type size() const noexcept { return m_size; }

	size_type max_size() const noexcept { return std::numeric_limits< size_type >::max() / sizeof(T); }

	bool empty() const noexcept { return m_size == 0; }

  private:
	T* m_data;
	size_type m_size;
	size_type m_dimensions;
	arrN m_dim_sz;
	arrN m_dim_steps;

	template< typename P >
	P at_inner(std::initializer_list< size_type > indexes) const
	{
		if (indexes.size() != N)
			throw std::invalid_argument("Demo_TensorLikeObject::at - incorrect number of indexes");
		size_t cur_ind = 0;
		auto it = indexes.begin();
		for (size_t i = 0; i < N; i++, ++it)
		{
			if (it == indexes.end())
				throw std::invalid_argument("Demo_TensorLikeObject::at - internal error");
			size_t ind = *it;
			if (ind >= m_dim_sz[i])
				throw std::out_of_range("Demo_TensorLikeObject::at");
			cur_ind += ind * m_dim_steps[i];
		}
		return m_data[cur_ind];
	}

	template< size_t F, typename Inits >
	void form_init_data(T* data, const Inits& inits, size_t& cur_ind)
	{
		if constexpr (F < N)
		{
			for (const auto& cur_init : inits)
				form_init_data< F + 1 >(data, cur_init, cur_ind);
		}
		else
		{
			new (&data[cur_ind++]) T(inits);
		}
	}

	void check_sz(arrN& sizes)
	{
		if constexpr (N == 0)
		{
			m_size = 1;
		}
		else
		{
			m_size = 1;
			for (int i = 0; i < N; i++)
			{
				m_dim_sz[i] = sizes[i];
				m_size *= sizes[i];
			}
			m_dim_steps[N - 1] = 1;
			for (int i = N - 2; i >= 0; i--)
				m_dim_steps[i] = m_dim_steps[i + 1] * m_dim_sz[i + 1];
		}
	}

	template< size_t F, typename I >
	void count_sz(arrN& sizes, const I& inits)
	{
		if constexpr (F < N)
		{
			sizes[F] = inits.size();
			if (!std::empty(inits))
				count_sz< F + 1 >(sizes, *(std::begin(inits)));
		}
	}

	void pre_fill_inits(arrN& sizes)
	{
		for (size_t i = 0; i < N; i++)
		{
			m_dim_sz[i] = sizes[i];
			m_size *= sizes[i];
		}
		if (N > 0)
		{
			m_dim_steps[N - 1] = 1;
			for (int i = N - 2; i >= 0; i--)
				m_dim_steps[i] = m_dim_steps[i + 1] * m_dim_sz[i + 1];
		}
	}

	template< typename V >
	void fill_data_with_inits_inner(size_t sz, V cur_val)
	{
		m_data = nullptr;
		if (sz == 0)
			return;
		void* cur_data = operator new(sz * sizeof(T));
		m_data = static_cast< T* >(cur_data);
		size_t cur_init_i = 0;
		try
		{
			for (; cur_init_i < sz; cur_init_i++)
			{
				new (&m_data[cur_init_i]) T(cur_val());
			}
		} catch (...)
		{
			act_erase(cur_init_i);
			operator delete(cur_data);
			throw;
		}
	}

	template< typename I >
	void fill_data_with_inits(size_t sz, I init_it)
	{
		fill_data_with_inits_inner(sz, [&init_it]() mutable -> T { return *init_it++; });
	}

	void fill_data_with_inits(size_t sz, const T& value)
	{
		fill_data_with_inits_inner(sz, [&value]() mutable -> T { return value; });
	}

	void act_erase(size_t stop_ind) noexcept
	{
		for (size_t i = stop_ind; i-- > 0;)
			m_data[i].~T();
	}
};
#endif
