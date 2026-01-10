<template>
  <div class="px-4 md:px-6 py-6 max-w-4xl mx-auto">
      <div class="flex justify-between items-center mb-6">
          <div class="flex items-center gap-3">
               <el-button icon="el-icon-back" circle @click="$router.push('/salev2/orders')"></el-button>
               <h1 class="text-2xl font-bold text-slate-800">Expense Manager</h1>
          </div>
          <div class="flex items-center gap-2">
               <el-date-picker v-model="filterMonth" type="month" placeholder="Pick a month" 
                    value-format="yyyy-MM-dd" @change="fetchExpenses" size="small" :clearable="false">
               </el-date-picker>
               <div class="bg-red-50 text-red-600 px-3 py-1 rounded font-bold border border-red-100">
                   Total: {{ formatPrice(summaryTotal) }}
               </div>
          </div>
      </div>

      <!-- Add New Expense Form -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5 mb-6">
          <h3 class="font-bold text-slate-700 mb-3 uppercase text-xs tracking-wider">Add New Expense</h3>
          <div class="flex flex-col md:flex-row gap-3 items-end">
              <div class="flex-1 w-full">
                  <label class="block text-xs font-bold text-slate-400 mb-1">Description</label>
                  <el-input v-model="form.description" placeholder="e.g., Import Goods, Repair..."></el-input>
              </div>
              <div class="w-full md:w-48">
                  <label class="block text-xs font-bold text-slate-400 mb-1">Amount</label>
                  <el-input-number v-model="form.amount" :step="100000" class="w-full" :min="0"></el-input-number>
              </div>
              <div class="w-full md:w-48">
                  <label class="block text-xs font-bold text-slate-400 mb-1">Date</label>
                  <el-date-picker v-model="form.entry_date" type="date" placeholder="Date" value-format="yyyy-MM-dd" class="w-full"></el-date-picker>
              </div>
              <el-button type="primary" icon="el-icon-plus" @click="addExpense" :loading="creating">Add</el-button>
          </div>
      </div>

      <!-- Expense List -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
          <el-table :data="expenses" v-loading="loading" style="width: 100%">
               <el-table-column label="Date" width="120">
                   <template slot-scope="scope">
                       <span class="font-medium text-slate-600">{{ formatDate(scope.row.entry_date) }}</span>
                   </template>
               </el-table-column>
               <el-table-column prop="description" label="Description">
                   <template slot-scope="scope">
                       <span class="font-bold text-slate-800">{{ scope.row.description }}</span>
                   </template>
               </el-table-column>
               <el-table-column label="Amount" width="150" align="right">
                   <template slot-scope="scope">
                       <span class="font-bold text-red-600">-{{ formatPrice(scope.row.amount) }}</span>
                   </template>
               </el-table-column>
               <el-table-column width="80" align="center">
                   <template slot-scope="scope">
                       <el-button type="text" icon="el-icon-delete" class="text-slate-400 hover:text-red-500" @click="deleteExpense(scope.row.id)"></el-button>
                   </template>
               </el-table-column>
          </el-table>
          
          <div class="p-4 flex justify-end" v-if="total > 20">
              <el-pagination background layout="prev, pager, next" :total="total" :page-size="20" @current-change="handlePageChange"></el-pagination>
          </div>
      </div>
  </div>
</template>

<script>
import moment from 'moment';

export default {
    data() {
        return {
            expenses: [],
            total: 0,
            summaryTotal: 0,
            loading: false,
            creating: false,
            filterMonth: moment().format('YYYY-MM-DD'),
            form: {
                description: '',
                amount: 0,
                entry_date: moment().format('YYYY-MM-DD')
            },
            page: 1
        }
    },
    methods: {
        formatPrice(value) {
            let val = (value/1).toFixed(0).replace('.', ',');
            return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".") + ' ₫';
        },
        formatDate(date) {
            return moment(date).format('DD/MM/YYYY');
        },
        async fetchExpenses() {
            this.loading = true;
            try {
                const res = await axios.get('/salev2/api/expenses', {
                    params: {
                        month: this.filterMonth,
                        page: this.page
                    }
                });
                this.expenses = res.data.expenses.data;
                this.total = res.data.expenses.total;
                // Assuming backend calculates TOTAL SUM of the filtered month, not just paginated page
                // We might need to adjust backend to return full sum
                 this.summaryTotal = res.data.expenses.data.reduce((sum, item) => sum + Number(item.amount), 0); 
                 // Wait, backend logic: query->sum('amount') was applied to the query... but paginate executes too. 
                 // Let's re-verify backend logic. It captures sum correctly.
                 this.summaryTotal = res.data.total;
            } catch (e) {
                console.error(e);
            } finally {
                this.loading = false;
            }
        },
        async addExpense() {
            if(!this.form.description || !this.form.amount) return this.$message.warning('Please fill info');
            this.creating = true;
            try {
                await axios.post('/salev2/api/expenses', this.form);
                this.$message.success('Expense Added');
                this.form.description = '';
                this.form.amount = 0;
                this.fetchExpenses();
            } catch (e) {
                this.$message.error('Failed');
            } finally {
                this.creating = false;
            }
        },
        async deleteExpense(id) {
            try {
                await this.$confirm('Are you sure?', 'Warning', { type: 'warning' });
                await axios.delete(`/salev2/api/expenses/${id}`);
                this.fetchExpenses();
            } catch(e) {}
        },
        handlePageChange(val) {
            this.page = val;
            this.fetchExpenses();
        }
    },
    mounted() {
        this.fetchExpenses();
    }
}
</script>
