<?php

namespace App\Http\Controllers;

use App\Expense;
use Illuminate\Http\Request;
use Carbon\Carbon;

class ExpenseController extends Controller
{
    public function index(Request $request) 
    {
        $query = Expense::orderBy('entry_date', 'desc')->orderBy('created_at', 'desc');

        if ($request->filled('month')) {
            $date = Carbon::parse($request->month);
            $query->whereYear('entry_date', $date->year)
                  ->whereMonth('entry_date', $date->month);
        }

        $expenses = $query->paginate(20);
        
        // Calculate total for filtered range (or all if not filtered, but usually filtered by month)
        $total = $query->sum('amount');

        return response()->json([
            'expenses' => $expenses,
            'total' => $total
        ]);
    }

    public function store(Request $request)
    {
        $request->validate([
            'description' => 'required|string',
            'amount' => 'required|numeric',
            'entry_date' => 'required|date'
        ]);

        $expense = Expense::create($request->all());

        return response()->json($expense);
    }

    public function destroy($id)
    {
        Expense::destroy($id);
        return response()->json(['success' => true]);
    }
}
